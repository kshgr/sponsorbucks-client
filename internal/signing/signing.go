package signing

import (
	"crypto/ed25519"
	"encoding/base64"
	"errors"
)

func SignBase64(privateKeyBase64 string, body []byte) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(privateKeyBase64)
	if err != nil {
		return "", err
	}
	if len(raw) != ed25519.PrivateKeySize {
		return "", errors.New("invalid ed25519 private key size")
	}
	sig := ed25519.Sign(ed25519.PrivateKey(raw), body)
	return base64.StdEncoding.EncodeToString(sig), nil
}
