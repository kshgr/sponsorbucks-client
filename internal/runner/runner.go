package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"sponsorbucks-client/internal/ads"
	"sponsorbucks-client/internal/api"
	"sponsorbucks-client/internal/buildinfo"
	"sponsorbucks-client/internal/events"
	"sponsorbucks-client/internal/idutil"
	"sponsorbucks-client/internal/localconfig"
	"sponsorbucks-client/internal/logs"
	"sponsorbucks-client/internal/overlay"
	"sponsorbucks-client/internal/signing"
	"sponsorbucks-client/internal/system"
	sbtools "sponsorbucks-client/internal/tools"
)

type Options struct {
	Surface string
	Command []string
	Version string
	Demo    bool
}

type adFrame struct {
	CampaignID string
	CreativeID string
	Line       string
}

var fallbackAds = []adFrame{
	{CampaignID: "camp_fallback_1", CreativeID: "creative_fallback_1", Line: "Sponsored - SponsorBucks launch partner"},
	{CampaignID: "camp_fallback_2", CreativeID: "creative_fallback_2", Line: "Sponsored - Keep shipping with SponsorBucks"},
	{CampaignID: "camp_fallback_3", CreativeID: "creative_fallback_3", Line: "Sponsored - Lightweight support for AI agents"},
}

func RunAgent(opts Options) (int, error) {
	if len(opts.Command) == 0 {
		return 1, fmt.Errorf("missing command")
	}

	cfg, err := localconfig.Load()
	if err != nil {
		return 1, err
	}

	surface := sbtools.CanonicalSurface(opts.Surface)
	activated := placementEnabled(cfg, surface)
	remoteReady := activated && !opts.Demo && cfg.APIBaseURL != "" && cfg.DeviceToken != "" && cfg.DeviceID != "" && cfg.DevicePrivateKey != ""
	client := api.New(cfg.APIBaseURL, cfg.DeviceToken)
	sessionID := idutil.NewUUID()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cmd := exec.CommandContext(ctx, opts.Command[0], opts.Command[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Start(); err != nil {
		return 1, err
	}

	var exited atomic.Bool
	var exitCode int32

	go func() {
		err := cmd.Wait()
		if err != nil {
			if ee, ok := err.(*exec.ExitError); ok {
				atomic.StoreInt32(&exitCode, int32(ee.ExitCode()))
			} else {
				atomic.StoreInt32(&exitCode, 1)
			}
		}
		exited.Store(true)
	}()

	if !activated {
		return waitForExit(ctx, &exited, &exitCode, cmd)
	}

	sticky := terminalSupportsSticky()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	currentAd := fallbackAd(0)
	if remoteReady {
		if ad, ok := fetchNextAd(client, cfg.DeviceID, surface, ""); ok {
			currentAd = ad
		}
		_ = postSigned(client, cfg, "/events-session-start", events.SessionStartEvent{
			EventType:      "session_start",
			SessionID:      sessionID,
			DeviceID:       cfg.DeviceID,
			Surface:        surface,
			StartedAt:      time.Now().UTC().Format(time.RFC3339),
			ClientVersion:  opts.Version,
			BuildID:        buildinfo.BuildID,
			BuildChannel:   buildinfo.BuildChannel,
			HumanInitiated: true,
		})
	}
	renderAd(sticky, currentAd.Line)

	sequence := 0
	lastCreative := currentAd.CreativeID

	for !exited.Load() {
		select {
		case <-ctx.Done():
			_ = cmd.Process.Signal(os.Interrupt)
		case <-ticker.C:
			if exited.Load() {
				continue
			}
			sequence++
			if remoteReady {
				if ad, ok := fetchNextAd(client, cfg.DeviceID, surface, lastCreative); ok {
					currentAd = ad
					lastCreative = ad.CreativeID
				}
			} else {
				currentAd = fallbackAd(sequence)
				lastCreative = currentAd.CreativeID
			}
			renderAd(sticky, currentAd.Line)

			if remoteReady {
				sig := system.CollectSignals()
				_ = postSigned(client, cfg, "/events-heartbeat", events.HeartbeatEvent{
					EventType:                  "heartbeat",
					SessionID:                  sessionID,
					DeviceID:                   cfg.DeviceID,
					Surface:                    surface,
					CampaignID:                 currentAd.CampaignID,
					CreativeID:                 currentAd.CreativeID,
					CreativeHash:               ads.CreativeHash(currentAd.Line),
					Sequence:                   sequence,
					VisibleMS:                  5000,
					ScreenUnlocked:             sig.ScreenUnlocked,
					RecentInputBucket:          sig.RecentInputBucket,
					ForegroundSupportedSurface: sig.ForegroundSupportedSurface,
					PlacementVisible:           true,
					CreatedAt:                  time.Now().UTC().Format(time.RFC3339),
					ClientVersion:              opts.Version,
					BuildID:                    buildinfo.BuildID,
					BuildChannel:               buildinfo.BuildChannel,
				})
			}
		}
	}

	if sticky {
		fmt.Fprintln(os.Stderr)
	}

	code := int(atomic.LoadInt32(&exitCode))
	if remoteReady {
		_ = postSigned(client, cfg, "/events-session-end", events.SessionEndEvent{
			EventType:     "session_end",
			SessionID:     sessionID,
			DeviceID:      cfg.DeviceID,
			Surface:       surface,
			EndedAt:       time.Now().UTC().Format(time.RFC3339),
			ExitCode:      code,
			ClientVersion: opts.Version,
			BuildID:       buildinfo.BuildID,
			BuildChannel:  buildinfo.BuildChannel,
		})
	}
	_ = logs.Append("run", map[string]string{
		"surface": surface,
		"mode":    map[bool]string{true: "demo", false: "live"}[opts.Demo],
		"exit":    fmt.Sprintf("%d", code),
	})
	return code, nil
}

func waitForExit(ctx context.Context, exited *atomic.Bool, exitCode *int32, cmd *exec.Cmd) (int, error) {
	for !exited.Load() {
		select {
		case <-ctx.Done():
			_ = cmd.Process.Signal(os.Interrupt)
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
	return int(atomic.LoadInt32(exitCode)), nil
}

func fetchNextAd(client *api.Client, deviceID, surface, lastCreative string) (adFrame, bool) {
	ad, err := client.NextAd(deviceID, surface, lastCreative)
	if err != nil || ad.NoAd || ad.Line == "" {
		return adFrame{}, false
	}
	return adFrame{
		CampaignID: ad.CampaignID,
		CreativeID: ad.CreativeID,
		Line:       ad.Line,
	}, true
}

func fallbackAd(sequence int) adFrame {
	if len(fallbackAds) == 0 {
		return adFrame{CampaignID: "camp_fallback", CreativeID: "creative_fallback", Line: "Sponsored - SponsorBucks"}
	}
	return fallbackAds[sequence%len(fallbackAds)]
}

func renderAd(sticky bool, line string) {
	if sticky {
		fmt.Fprintf(os.Stderr, "\r\033[2K[%s] %s", time.Now().Format("15:04:05"), overlay.SponsoredLine(line))
		return
	}
	fmt.Fprintf(os.Stderr, "\n[%s] %s\n", time.Now().Format("15:04:05"), overlay.SponsoredLine(line))
}

func postSigned(client *api.Client, cfg localconfig.Config, path string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	ts := time.Now().UTC().Format(time.RFC3339)
	sig, err := signing.SignBase64(cfg.DevicePrivateKey, body)
	if err != nil {
		return err
	}
	return client.PostSigned(path, body, cfg.DeviceID, sig, ts)
}

func terminalSupportsSticky() bool {
	if fi, err := os.Stderr.Stat(); err == nil && (fi.Mode()&os.ModeCharDevice) != 0 {
		return true
	}
	if fi, err := os.Stdout.Stat(); err == nil && (fi.Mode()&os.ModeCharDevice) != 0 {
		return true
	}
	return false
}

func placementEnabled(cfg localconfig.Config, surface string) bool {
	return !cfg.Paused && !cfg.DisabledTools[surface]
}
