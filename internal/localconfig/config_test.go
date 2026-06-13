package localconfig

import "testing"

func TestConfigRedaction(t *testing.T) {
	cfg := Config{
		APIBaseURL:       "https://example.invalid",
		DeviceID:         "dev_123",
		UserID:           "user_123",
		DeviceToken:      "token_secret",
		DevicePublicKey:  "public_key",
		DevicePrivateKey: "private_secret",
		LinkCode:         "link_secret",
	}
	redacted := cfg.Redacted()
	if redacted.DeviceToken != "***redacted***" {
		t.Fatalf("device token not redacted: %q", redacted.DeviceToken)
	}
	if redacted.DevicePrivateKey != "***redacted***" {
		t.Fatalf("private key not redacted: %q", redacted.DevicePrivateKey)
	}
	if redacted.LinkCode != "***redacted***" {
		t.Fatalf("link code not redacted: %q", redacted.LinkCode)
	}
	if redacted.DeviceID != cfg.DeviceID || redacted.UserID != cfg.UserID || redacted.APIBaseURL != cfg.APIBaseURL {
		t.Fatalf("non-secret fields changed during redaction")
	}
}
