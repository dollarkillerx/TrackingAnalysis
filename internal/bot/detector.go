package bot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tracking/analysis/internal/config"
)

var botPatterns = []string{
	"bot", "spider", "crawler", "headless", "phantom",
	"selenium", "puppeteer", "scrapy", "wget", "curl",
}

func Score(ua, acceptLang, secFetch, referer string, recentHits int) int {
	score := 0
	uaLower := strings.ToLower(ua)

	// Check for bot-like user agents
	for _, pattern := range botPatterns {
		if strings.Contains(uaLower, pattern) {
			score += 50
			break
		}
	}

	// No Accept-Language header
	if acceptLang == "" {
		score += 20
	}

	// All Sec-Fetch-* headers missing
	if secFetch == "" {
		score += 20
	}

	// Suspicious frequency with no referer
	if referer == "" && recentHits > 10 {
		score += 30
	}

	if score > 100 {
		score = 100
	}
	return score
}

func IsBot(score int, cfg *config.BotConfiguration) (blocked bool, suspected bool) {
	if score >= cfg.BlockThreshold {
		return true, true
	}
	if score >= cfg.MarkThreshold {
		return false, true
	}
	return false, false
}

func CountRecentHits(ctx context.Context, rdb *redis.Client, ip string) int {
	key := fmt.Sprintf("bot:freq:%s", ip)
	val, err := rdb.Incr(ctx, key).Result()
	if err != nil {
		return 0
	}
	if val == 1 {
		rdb.Expire(ctx, key, 10*time.Second)
	}
	return int(val)
}
