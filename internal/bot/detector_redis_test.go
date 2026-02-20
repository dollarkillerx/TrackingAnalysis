package bot

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func setupRedis(t *testing.T) (*miniredis.Miniredis, *redis.Client) {
	t.Helper()
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	return mr, rdb
}

func TestCountRecentHits_Increment(t *testing.T) {
	_, rdb := setupRedis(t)
	ctx := context.Background()

	count1 := CountRecentHits(ctx, rdb, "192.168.1.1")
	if count1 != 1 {
		t.Errorf("first call = %d, want 1", count1)
	}

	count2 := CountRecentHits(ctx, rdb, "192.168.1.1")
	if count2 != 2 {
		t.Errorf("second call = %d, want 2", count2)
	}

	count3 := CountRecentHits(ctx, rdb, "192.168.1.1")
	if count3 != 3 {
		t.Errorf("third call = %d, want 3", count3)
	}
}

func TestCountRecentHits_DifferentIPs(t *testing.T) {
	_, rdb := setupRedis(t)
	ctx := context.Background()

	CountRecentHits(ctx, rdb, "10.0.0.1")
	CountRecentHits(ctx, rdb, "10.0.0.1")
	CountRecentHits(ctx, rdb, "10.0.0.1")

	count := CountRecentHits(ctx, rdb, "10.0.0.2") // different IP
	if count != 1 {
		t.Errorf("different IP first call = %d, want 1", count)
	}
}
