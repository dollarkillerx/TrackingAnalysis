package security

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	payload := TokenPayload{
		TrackerID: "tracker-1",
		TargetID:  "target-1",
		Mode:      "js",
		Exp:       time.Now().Add(time.Hour).Unix(),
	}
	token, err := GenerateToken(payload, "secret")
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}
	parts := strings.SplitN(token, ".", 2)
	if len(parts) != 2 {
		t.Fatalf("expected token format base64.base64, got %q", token)
	}
	if parts[0] == "" || parts[1] == "" {
		t.Fatal("token parts must not be empty")
	}
}

func TestVerifyToken_Valid(t *testing.T) {
	secret := "test-secret"
	payload := TokenPayload{
		TrackerID:  "tracker-1",
		CampaignID: "campaign-1",
		ChannelID:  "channel-1",
		TargetID:   "target-1",
		Exp:        time.Now().Add(time.Hour).Unix(),
		Mode:       "302",
	}
	token, err := GenerateToken(payload, secret)
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}
	got, err := VerifyToken(token, secret)
	if err != nil {
		t.Fatalf("VerifyToken: %v", err)
	}
	if got.TrackerID != payload.TrackerID {
		t.Errorf("TrackerID = %q, want %q", got.TrackerID, payload.TrackerID)
	}
	if got.CampaignID != payload.CampaignID {
		t.Errorf("CampaignID = %q, want %q", got.CampaignID, payload.CampaignID)
	}
	if got.Mode != payload.Mode {
		t.Errorf("Mode = %q, want %q", got.Mode, payload.Mode)
	}
}

func TestVerifyToken_InvalidSignature(t *testing.T) {
	token, _ := GenerateToken(TokenPayload{TrackerID: "t1", TargetID: "t2"}, "secret1")
	_, err := VerifyToken(token, "wrong-secret")
	if err == nil {
		t.Fatal("expected error for invalid signature")
	}
	if !strings.Contains(err.Error(), "signature") {
		t.Errorf("error = %q, want mention of signature", err.Error())
	}
}

func TestVerifyToken_InvalidFormat(t *testing.T) {
	_, err := VerifyToken("no-dot-separator", "secret")
	if err == nil {
		t.Fatal("expected error for missing dot separator")
	}
	if !strings.Contains(err.Error(), "format") {
		t.Errorf("error = %q, want mention of format", err.Error())
	}
}

func TestVerifyToken_Expired(t *testing.T) {
	payload := TokenPayload{
		TrackerID: "t1",
		TargetID:  "t2",
		Exp:       time.Now().Add(-time.Hour).Unix(), // expired 1 hour ago
	}
	token, _ := GenerateToken(payload, "secret")
	_, err := VerifyToken(token, "secret")
	if err == nil {
		t.Fatal("expected error for expired token")
	}
	if !strings.Contains(err.Error(), "expired") {
		t.Errorf("error = %q, want mention of expired", err.Error())
	}
}

func TestVerifyToken_NoExpiry(t *testing.T) {
	payload := TokenPayload{
		TrackerID: "t1",
		TargetID:  "t2",
		Exp:       0, // no expiration
	}
	token, _ := GenerateToken(payload, "secret")
	got, err := VerifyToken(token, "secret")
	if err != nil {
		t.Fatalf("VerifyToken should succeed for Exp=0: %v", err)
	}
	if got.TrackerID != "t1" {
		t.Errorf("TrackerID = %q, want %q", got.TrackerID, "t1")
	}
}
