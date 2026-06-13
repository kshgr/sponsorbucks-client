package commands

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"sponsorbucks-client/internal/api"
	"sponsorbucks-client/internal/buildinfo"
	"sponsorbucks-client/internal/device"
	"sponsorbucks-client/internal/localconfig"
	"sponsorbucks-client/internal/logs"
	"sponsorbucks-client/internal/openurl"
)

func Login(args []string, info buildinfo.Info) {
	cfg, err := localconfig.Load()
	exitOnErr(err)

	if cfg.APIBaseURL == "" {
		cfg.APIBaseURL = api.DefaultBaseURL
		exitOnErr(localconfig.Save(cfg))
		fmt.Printf("Using default SponsorBucks API: %s\n", cfg.APIBaseURL)
	}

	if cfg.DevicePrivateKey == "" || cfg.DevicePublicKey == "" {
		keys, err := device.GenerateKeyPair()
		exitOnErr(err)
		cfg.DevicePrivateKey = keys.PrivateKeyBase64
		cfg.DevicePublicKey = keys.PublicKeyBase64
		exitOnErr(localconfig.Save(cfg))
	}

	client := api.New(cfg.APIBaseURL, "")
	resp, err := client.StartLink(api.StartLinkRequest{
		DevicePublicKey: cfg.DevicePublicKey,
		DeviceName:      device.DefaultDeviceName(),
		ClientVersion:   info.ClientVersion,
		BuildID:         info.BuildID,
		BuildChannel:    info.BuildChannel,
		OS:              runtime.GOOS,
		Arch:            runtime.GOARCH,
	})
	exitOnErr(err)

	cfg.DeviceID = resp.DeviceID
	cfg.LinkCode = resp.LinkCode
	exitOnErr(localconfig.Save(cfg))

	fmt.Println("Opening browser to link this device:")
	fmt.Println(resp.LinkURL)
	_ = openurl.Open(resp.LinkURL)

	fmt.Println("Waiting for browser sign-in...")
	for i := 0; i < 120; i++ {
		linked, err := client.CompleteLink(api.CompleteLinkRequest{
			DeviceID: resp.DeviceID,
			LinkCode: resp.LinkCode,
		})
		if err == nil && linked.Status == "linked" {
			cfg.DeviceToken = linked.DeviceToken
			cfg.UserID = linked.UserID
			cfg.LinkCode = ""
			exitOnErr(localconfig.Save(cfg))
			fmt.Println("Device linked successfully.")
			_ = logs.Append("login", nil)
			return
		}
		time.Sleep(2 * time.Second)
	}

	fmt.Println("Timed out waiting for device link. Run sponsorbucks login again.")
	os.Exit(1)
}
