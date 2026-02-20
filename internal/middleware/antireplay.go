package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tracking/analysis/internal/config"
)

func CheckAntiReplay(ctx context.Context, rdb *redis.Client, cfg *config.SecurityConfiguration, ts int64, nonce2 string) error {
	now := time.Now().Unix()
	diff := now - ts
	if diff < 0 {
		diff = -diff
	}
	if diff > int64(cfg.TSWindowSeconds) {
		return fmt.Errorf("expired_timestamp")
	}

	// Nonce replay check
	key := fmt.Sprintf("nonce:%s", nonce2)
	ttl := time.Duration(cfg.NonceTTLSeconds) * time.Second
	set, err := rdb.SetNX(ctx, key, 1, ttl).Result()
	if err != nil {
		return nil // fail open on Redis error
	}
	if !set {
		return fmt.Errorf("replay_detected")
	}

	return nil
}

