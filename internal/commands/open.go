package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"sponsorbucks-client/internal/adstate"
	"sponsorbucks-client/internal/api"
	"sponsorbucks-client/internal/buildinfo"
	"sponsorbucks-client/internal/events"
	"sponsorbucks-client/internal/localconfig"
	"sponsorbucks-client/internal/logs"
	"sponsorbucks-client/internal/openurl"
	"sponsorbucks-client/internal/signing"
)

func Open(args []string) {
	_ = args
	state, err := adstate.Load()
	exitOnErr(err)
	if !state.Visible || state.OpenURL == "" {
		fmt.Println("No visible sponsored ad is cached right now.")
		os.Exit(1)
	}

	if err := openurl.Open(state.OpenURL); err != nil {
		exitOnErr(err)
	}

	_ = logs.Append("click", map[string]string{
		"campaign_id": state.CampaignID,
		"creative_id": state.CreativeID,
		"surface":     state.Surface,
	})

	cfg, err := localconfig.Load()
	exitOnErr(err)
	if cfg.DeviceToken == "" || cfg.DeviceID == "" || cfg.DevicePrivateKey == "" || cfg.APIBaseURL == "" {
		fmt.Println("Opened cached ad URL locally.")
		return
	}

	client := api.New(cfg.APIBaseURL, cfg.DeviceToken)
	body, err := jsonMarshal(events.ClickEvent{
		EventType:     "click",
		SessionID:     state.SessionID,
		DeviceID:      cfg.DeviceID,
		CampaignID:    state.CampaignID,
		CreativeID:    state.CreativeID,
		CreativeHash:  state.CreativeHash,
		Surface:       state.Surface,
		ClickedAt:     time.Now().UTC().Format(time.RFC3339),
		ClientVersion: buildinfo.Version,
		BuildID:       buildinfo.BuildID,
		BuildChannel:  buildinfo.BuildChannel,
	})
	exitOnErr(err)

	ts := time.Now().UTC().Format(time.RFC3339)
	sig, err := signing.SignBase64(cfg.DevicePrivateKey, body)
	exitOnErr(err)
	if err := client.PostSigned("/events-click", body, cfg.DeviceID, sig, ts); err != nil {
		exitOnErr(err)
	}

	fmt.Println("Opened cached ad URL and logged click.")
}

func jsonMarshal(v any) ([]byte, error) {
	return json.Marshal(v)
}
