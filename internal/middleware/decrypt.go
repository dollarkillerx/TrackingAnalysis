package middleware

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"

	"github.com/tracking/analysis/internal/security"
)

type EncryptedParams struct {
	EK    string `json:"ek"`
	Nonce string `json:"nonce"`
	CT    string `json:"ct"`
	TS    int64  `json:"ts"`
	Nonce2 string `json:"nonce2"`
	KID   string `json:"kid"`
}

func DecryptRequest(params EncryptedParams, privKey *rsa.PrivateKey) ([]byte, error) {
	ek, err := base64.StdEncoding.DecodeString(params.EK)
	if err != nil {
		return nil, fmt.Errorf("invalid ek encoding")
	}

	nonce, err := base64.StdEncoding.DecodeString(params.Nonce)
	if err != nil {
		return nil, fmt.Errorf("invalid nonce encoding")
	}

	ct, err := base64.StdEncoding.DecodeString(params.CT)
	if err != nil {
		return nil, fmt.Errorf("invalid ct encoding")
	}

	plaintext, err := security.DecryptPayload(ek, nonce, ct, privKey)
	if err != nil {
		return nil, fmt.Errorf("decrypt_failed: %w", err)
	}

	return plaintext, nil
}
