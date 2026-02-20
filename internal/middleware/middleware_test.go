package middleware

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"strings"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/tracking/analysis/internal/config"
)

func setupRedis(t *testing.T) (*miniredis.Miniredis, *redis.Client) {
	t.Helper()
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	return mr, rdb
}

// --- Rate Limit Tests ---

func TestCheckRateLimit_UnderLimit(t *testing.T) {
	_, rdb := setupRedis(t)
	cfg := &config.RateLimitConfiguration{
		PerIPPerMinute:        10,
		PerIPUAPerMinute:      10,
		PerTrackerIPPerMinute: 10,
	}
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		if err := CheckRateLimit(ctx, rdb, cfg, "1.2.3.4", "Mozilla/5.0", "tracker-1"); err != nil {
			t.Fatalf("request %d: unexpected error: %v", i, err)
		}
	}
}

func TestCheckRateLimit_OverLimit(t *testing.T) {
	_, rdb := setupRedis(t)
	cfg := &config.RateLimitConfiguration{
		PerIPPerMinute:        3,
		PerIPUAPerMinute:      100,
		PerTrackerIPPerMinute: 100,
	}
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		if err := CheckRateLimit(ctx, rdb, cfg, "1.2.3.4", "Mozilla/5.0", ""); err != nil {
			t.Fatalf("request %d: unexpected error: %v", i, err)
		}
	}

	err := CheckRateLimit(ctx, rdb, cfg, "1.2.3.4", "Mozilla/5.0", "")
	if err == nil {
		t.Fatal("expected rate_limited error")
	}
	if !strings.Contains(err.Error(), "rate_limited") {
		t.Errorf("error = %q, want rate_limited", err.Error())
	}
}

func TestCheckRateLimit_PerIPUALimit(t *testing.T) {
	_, rdb := setupRedis(t)
	cfg := &config.RateLimitConfiguration{
		PerIPPerMinute:        100,
		PerIPUAPerMinute:      2,
		PerTrackerIPPerMinute: 100,
	}
	ctx := context.Background()

	for i := 0; i < 2; i++ {
		if err := CheckRateLimit(ctx, rdb, cfg, "1.2.3.4", "SpecificUA", ""); err != nil {
			t.Fatalf("request %d: unexpected error: %v", i, err)
		}
	}

	err := CheckRateLimit(ctx, rdb, cfg, "1.2.3.4", "SpecificUA", "")
	if err == nil {
		t.Fatal("expected rate_limited error for IP+UA limit")
	}
}

func TestCheckRateLimit_TrackerIPLimit(t *testing.T) {
	_, rdb := setupRedis(t)
	cfg := &config.RateLimitConfiguration{
		PerIPPerMinute:        100,
		PerIPUAPerMinute:      100,
		PerTrackerIPPerMinute: 2,
	}
	ctx := context.Background()

	for i := 0; i < 2; i++ {
		if err := CheckRateLimit(ctx, rdb, cfg, "1.2.3.4", "Mozilla", "tracker-1"); err != nil {
			t.Fatalf("request %d: unexpected error: %v", i, err)
		}
	}

	err := CheckRateLimit(ctx, rdb, cfg, "1.2.3.4", "Mozilla", "tracker-1")
	if err == nil {
		t.Fatal("expected rate_limited error for tracker+IP limit")
	}
}

// --- Anti-Replay Tests ---

func TestCheckAntiReplay_ValidRequest(t *testing.T) {
	_, rdb := setupRedis(t)
	cfg := &config.SecurityConfiguration{
		TSWindowSeconds: 30,
		NonceTTLSeconds: 60,
	}
	ctx := context.Background()

	ts := time.Now().Unix()
	err := CheckAntiReplay(ctx, rdb, cfg, ts, "unique-nonce-1")
	if err != nil {
		t.Fatalf("valid request returned error: %v", err)
	}
}

func TestCheckAntiReplay_ExpiredTimestamp(t *testing.T) {
	_, rdb := setupRedis(t)
	cfg := &config.SecurityConfiguration{
		TSWindowSeconds: 30,
		NonceTTLSeconds: 60,
	}
	ctx := context.Background()

	ts := time.Now().Add(-time.Minute).Unix() // 60 seconds ago, window is 30
	err := CheckAntiReplay(ctx, rdb, cfg, ts, "nonce-expired")
	if err == nil {
		t.Fatal("expected expired_timestamp error")
	}
	if !strings.Contains(err.Error(), "expired_timestamp") {
		t.Errorf("error = %q, want expired_timestamp", err.Error())
	}
}

func TestCheckAntiReplay_ReplayedNonce(t *testing.T) {
	_, rdb := setupRedis(t)
	cfg := &config.SecurityConfiguration{
		TSWindowSeconds: 30,
		NonceTTLSeconds: 60,
	}
	ctx := context.Background()
	ts := time.Now().Unix()

	// First request — should succeed
	if err := CheckAntiReplay(ctx, rdb, cfg, ts, "nonce-replay"); err != nil {
		t.Fatalf("first request: %v", err)
	}

	// Replayed nonce — should fail
	err := CheckAntiReplay(ctx, rdb, cfg, ts, "nonce-replay")
	if err == nil {
		t.Fatal("expected replay_detected error")
	}
	if !strings.Contains(err.Error(), "replay_detected") {
		t.Errorf("error = %q, want replay_detected", err.Error())
	}
}

// --- Decrypt Tests ---

func generateTestRSAKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("rsa.GenerateKey: %v", err)
	}
	return key
}

func hybridEncrypt(t *testing.T, plaintext []byte, pub *rsa.PublicKey) (ekB64, nonceB64, ctB64 string) {
	t.Helper()
	aesKey := make([]byte, 32)
	if _, err := rand.Read(aesKey); err != nil {
		t.Fatalf("rand.Read: %v", err)
	}
	ek, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, aesKey, nil)
	if err != nil {
		t.Fatalf("EncryptOAEP: %v", err)
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		t.Fatalf("NewCipher: %v", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		t.Fatalf("NewGCM: %v", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		t.Fatalf("rand.Read nonce: %v", err)
	}
	ct := gcm.Seal(nil, nonce, plaintext, nil)

	return base64.StdEncoding.EncodeToString(ek),
		base64.StdEncoding.EncodeToString(nonce),
		base64.StdEncoding.EncodeToString(ct)
}

func TestDecryptRequest_Roundtrip(t *testing.T) {
	key := generateTestRSAKey(t)
	plaintext := []byte(`{"tracker_id":"abc","visitor_id":"xyz"}`)

	ekB64, nonceB64, ctB64 := hybridEncrypt(t, plaintext, &key.PublicKey)

	params := EncryptedParams{
		EK:    ekB64,
		Nonce: nonceB64,
		CT:    ctB64,
	}

	got, err := DecryptRequest(params, key)
	if err != nil {
		t.Fatalf("DecryptRequest: %v", err)
	}
	if string(got) != string(plaintext) {
		t.Errorf("decrypted = %q, want %q", got, plaintext)
	}
}

func TestDecryptRequest_InvalidBase64(t *testing.T) {
	key := generateTestRSAKey(t)

	tests := []struct {
		name   string
		params EncryptedParams
	}{
		{
			name: "bad ek",
			params: EncryptedParams{
				EK:    "!!!not-base64!!!",
				Nonce: base64.StdEncoding.EncodeToString([]byte("nonce")),
				CT:    base64.StdEncoding.EncodeToString([]byte("ct")),
			},
		},
		{
			name: "bad nonce",
			params: EncryptedParams{
				EK:    base64.StdEncoding.EncodeToString([]byte("ek")),
				Nonce: "!!!not-base64!!!",
				CT:    base64.StdEncoding.EncodeToString([]byte("ct")),
			},
		},
		{
			name: "bad ct",
			params: EncryptedParams{
				EK:    base64.StdEncoding.EncodeToString([]byte("ek")),
				Nonce: base64.StdEncoding.EncodeToString([]byte("nonce")),
				CT:    "!!!not-base64!!!",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecryptRequest(tt.params, key)
			if err == nil {
				t.Fatalf("expected error for %s", tt.name)
			}
		})
	}
}
