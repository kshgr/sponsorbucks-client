package device

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"os"
)

type KeyPair struct {
	PublicKeyBase64  string
	PrivateKeyBase64 string
}

func GenerateKeyPair() (KeyPair, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return KeyPair{}, err
	}
	return KeyPair{
		PublicKeyBase64:  base64.StdEncoding.EncodeToString(pub),
		PrivateKeyBase64: base64.StdEncoding.EncodeToString(priv),
	}, nil
}

func DefaultDeviceName() string {
	host, err := os.Hostname()
	if err != nil || host == "" {
		return "SponsorBucks device"
	}
	return host
}
