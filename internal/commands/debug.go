package commands

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"sponsorbucks-client/internal/ads"
	"sponsorbucks-client/internal/buildinfo"
	"sponsorbucks-client/internal/events"
	"sponsorbucks-client/internal/idutil"
	"sponsorbucks-client/internal/localconfig"
	"sponsorbucks-client/internal/overlay"
)

func Debug(args []string, info buildinfo.Info) {
	if len(args) == 0 || args[0] != "event" {
		fmt.Println("Usage: sponsorbucks debug event [--json|--explain]")
		os.Exit(1)
	}

	fs := flag.NewFlagSet("event", flag.ExitOnError)
	jsonOnly := fs.Bool("json", false, "print JSON only")
	explain := fs.Bool("explain", false, "print JSON and field explanations")
	_ = fs.Parse(args[1:])

	cfg, err := localconfig.Load()
	exitOnErr(err)

	visibleLine := overlay.SponsoredLine("Sponsored - demo creative")
	ev := events.HeartbeatEvent{
		EventType:                  "heartbeat",
		SessionID:                  idutil.NewUUID(),
		DeviceID:                   cfg.DeviceID,
		Surface:                    "codex",
		CampaignID:                 "camp_demo",
		CreativeID:                 "creative_demo",
		CreativeHash:               ads.CreativeHash(visibleLine),
		Sequence:                   1,
		VisibleMS:                  5000,
		ScreenUnlocked:             true,
		RecentInputBucket:          "under_2_minutes",
		ForegroundSupportedSurface: true,
		PlacementVisible:           true,
		CreatedAt:                  time.Now().UTC().Format(time.RFC3339),
		ClientVersion:              info.ClientVersion,
		BuildID:                    info.BuildID,
		BuildChannel:               info.BuildChannel,
	}

	out, err := json.MarshalIndent(ev, "", "  ")
	exitOnErr(err)
	if *jsonOnly {
		fmt.Println(string(out))
		return
	}

	fmt.Println(string(out))
	if !*explain {
		return
	}

	fmt.Println()
	fmt.Println("Field meanings:")
	fmt.Println("- event_type: event kind")
	fmt.Println("- session_id: local session identifier")
	fmt.Println("- device_id: linked device identifier")
	fmt.Println("- surface: supported tool surface")
	fmt.Println("- campaign_id: selected campaign")
	fmt.Println("- creative_id: selected creative")
	fmt.Println("- creative_hash: SHA-256 of the visible sponsored line")
	fmt.Println("- sequence: heartbeat number within this session")
	fmt.Println("- visible_ms: intended visible duration")
	fmt.Println("- screen_unlocked: best-effort lock state")
	fmt.Println("- recent_input_bucket: coarse recent activity bucket")
	fmt.Println("- foreground_supported_surface: wrapper/adapter visibility flag")
	fmt.Println("- placement_visible: whether the placement was visible")
	fmt.Println("- created_at: event timestamp")
	fmt.Println("- client_version: client version string")
	fmt.Println("- build_id: build identifier")
	fmt.Println("- build_channel: dev or release channel")
	fmt.Printf("- visible line: %s\n", visibleLine)
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
