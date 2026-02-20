package middleware

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tracking/analysis/internal/config"
)

func CheckRateLimit(ctx context.Context, rdb *redis.Client, cfg *config.RateLimitConfiguration, ip, ua, trackerID string) error {
	now := time.Now().Unix()
	window := now / 60 // 1-minute window

	// Per-IP limit
	ipKey := fmt.Sprintf("rl:ip:%s:%d", ip, window)
	if err := checkLimit(ctx, rdb, ipKey, cfg.PerIPPerMinute); err != nil {
		return err
	}

	// Per-IP+UA limit
	uaHash := fmt.Sprintf("%x", sha256.Sum256([]byte(ua)))[:16]
	ipuaKey := fmt.Sprintf("rl:ipua:%s:%s:%d", ip, uaHash, window)
	if err := checkLimit(ctx, rdb, ipuaKey, cfg.PerIPUAPerMinute); err != nil {
		return err
	}

	// Per-Tracker+IP limit (only if tracker specified)
	if trackerID != "" {
		tKey := fmt.Sprintf("rl:tracker_ip:%s:%s:%d", trackerID, ip, window)
		if err := checkLimit(ctx, rdb, tKey, cfg.PerTrackerIPPerMinute); err != nil {
			return err
		}
	}

	return nil
}

func checkLimit(ctx context.Context, rdb *redis.Client, key string, limit int) error {
	val, err := rdb.Incr(ctx, key).Result()
	if err != nil {
		return nil // fail open on Redis error
	}
	if val == 1 {
		rdb.Expire(ctx, key, 60*time.Second)
	}
	if int(val) > limit {
		return fmt.Errorf("rate_limited")
	}
	return nil
}
