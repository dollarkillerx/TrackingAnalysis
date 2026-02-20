package dedup

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tracking/analysis/internal/config"
)

func CheckClickDedup(ctx context.Context, rdb *redis.Client, cfg *config.SecurityConfiguration, trackerID, channelID, visitorID string) bool {
	key := fmt.Sprintf("dedup:click:%s:%s:%s", trackerID, channelID, visitorID)
	ttl := time.Duration(cfg.DedupSeconds) * time.Second
	set, err := rdb.SetNX(ctx, key, 1, ttl).Result()
	if err != nil {
		return false // fail open
	}
	return !set // isDuplicate = true if key already existed
}
