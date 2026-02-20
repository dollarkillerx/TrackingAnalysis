package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type TokenPayload struct {
	TrackerID  string `json:"tracker_id"`
	CampaignID string `json:"campaign_id,omitempty"`
	ChannelID  string `json:"channel_id,omitempty"`
	TargetID   string `json:"target_id"`
	Exp        int64  `json:"exp"`
	Mode       string `json:"mode"` // "js" or "302"
}

func GenerateToken(payload TokenPayload, secret string) (string, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	encoded := base64.RawURLEncoding.EncodeToString(data)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(encoded))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return encoded + "." + sig, nil
}

func VerifyToken(token, secret string) (*TokenPayload, error) {
	parts := strings.SplitN(token, ".", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid token format")
	}
	encoded, sig := parts[0], parts[1]

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(encoded))
	expectedSig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(sig), []byte(expectedSig)) {
		return nil, errors.New("invalid token signature")
	}

	data, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, errors.New("invalid token encoding")
	}

	var payload TokenPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, errors.New("invalid token payload")
	}

	if payload.Exp > 0 && time.Now().Unix() > payload.Exp {
		return nil, errors.New("token expired")
	}

	return &payload, nil
}
