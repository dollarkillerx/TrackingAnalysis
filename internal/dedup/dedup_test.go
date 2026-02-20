package dedup

import (
	"context"
	"testing"

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

func TestCheckClickDedup_FirstClick(t *testing.T) {
	_, rdb := setupRedis(t)
	cfg := &config.SecurityConfiguration{DedupSeconds: 3600}
	ctx := context.Background()

	isDup := CheckClickDedup(ctx, rdb, cfg, "tracker-1", "channel-1", "visitor-1")
	if isDup {
		t.Error("first click should not be duplicate")
	}
}

func TestCheckClickDedup_DuplicateClick(t *testing.T) {
	_, rdb := setupRedis(t)
	cfg := &config.SecurityConfiguration{DedupSeconds: 3600}
	ctx := context.Background()

	CheckClickDedup(ctx, rdb, cfg, "tracker-1", "channel-1", "visitor-1")
	isDup := CheckClickDedup(ctx, rdb, cfg, "tracker-1", "channel-1", "visitor-1")
	if !isDup {
		t.Error("second click with same IDs should be duplicate")
	}
}

func TestCheckClickDedup_DifferentVisitor(t *testing.T) {
	_, rdb := setupRedis(t)
	cfg := &config.SecurityConfiguration{DedupSeconds: 3600}
	ctx := context.Background()

	CheckClickDedup(ctx, rdb, cfg, "tracker-1", "channel-1", "visitor-1")
	isDup := CheckClickDedup(ctx, rdb, cfg, "tracker-1", "channel-1", "visitor-2")
	if isDup {
		t.Error("different visitor_id should not be duplicate")
	}
}

func TestCheckClickDedup_TTLExpiry(t *testing.T) {
	mr, rdb := setupRedis(t)
	cfg := &config.SecurityConfiguration{DedupSeconds: 10}
	ctx := context.Background()

	CheckClickDedup(ctx, rdb, cfg, "tracker-1", "channel-1", "visitor-1")

	// Fast-forward miniredis past the TTL
	mr.FastForward(11 * 1e9) // 11 seconds in nanoseconds

	isDup := CheckClickDedup(ctx, rdb, cfg, "tracker-1", "channel-1", "visitor-1")
	if isDup {
		t.Error("after TTL expiry, same click should not be duplicate")
	}
}
