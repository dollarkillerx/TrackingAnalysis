package bot

import (
	"testing"

	"github.com/tracking/analysis/internal/config"
)

func TestScore_NormalBrowser(t *testing.T) {
	score := Score(
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		"en-US,en;q=0.9",
		"navigate",
		"https://example.com",
		0,
	)
	if score != 0 {
		t.Errorf("normal browser score = %d, want 0", score)
	}
}

func TestScore_BotUA(t *testing.T) {
	score := Score(
		"Googlebot/2.1 (+http://www.google.com/bot.html)",
		"en-US",
		"navigate",
		"https://google.com",
		0,
	)
	if score < 50 {
		t.Errorf("bot UA score = %d, want >= 50", score)
	}
}

func TestScore_HeadlessUA(t *testing.T) {
	score := Score(
		"Mozilla/5.0 HeadlessChrome/90.0",
		"en-US",
		"navigate",
		"https://example.com",
		0,
	)
	if score < 50 {
		t.Errorf("headless UA score = %d, want >= 50", score)
	}
}

func TestScore_NoAcceptLanguage(t *testing.T) {
	score := Score(
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		"",        // no Accept-Language
		"navigate",
		"https://example.com",
		0,
	)
	if score < 20 {
		t.Errorf("no Accept-Language score = %d, want >= 20", score)
	}
}

func TestScore_NoSecFetch(t *testing.T) {
	score := Score(
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		"en-US",
		"", // no Sec-Fetch
		"https://example.com",
		0,
	)
	if score < 20 {
		t.Errorf("no Sec-Fetch score = %d, want >= 20", score)
	}
}

func TestScore_SuspiciousFrequency(t *testing.T) {
	score := Score(
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		"en-US",
		"navigate",
		"", // no referer
		15, // > 10 recent hits
	)
	if score < 30 {
		t.Errorf("suspicious frequency score = %d, want >= 30", score)
	}
}

func TestScore_MaxCap(t *testing.T) {
	// Bot UA (+50) + no Accept-Language (+20) + no Sec-Fetch (+20) + suspicious freq (+30) = 120 → capped at 100
	score := Score(
		"Googlebot/2.1",
		"",
		"",
		"",
		15,
	)
	if score > 100 {
		t.Errorf("score = %d, want <= 100", score)
	}
	if score != 100 {
		t.Errorf("max combined score = %d, want 100", score)
	}
}

func TestIsBot_BelowThresholds(t *testing.T) {
	cfg := &config.BotConfiguration{MarkThreshold: 50, BlockThreshold: 80}
	blocked, suspected := IsBot(30, cfg)
	if blocked {
		t.Error("score 30 should not be blocked")
	}
	if suspected {
		t.Error("score 30 should not be suspected")
	}
}

func TestIsBot_SuspectedOnly(t *testing.T) {
	cfg := &config.BotConfiguration{MarkThreshold: 50, BlockThreshold: 80}
	blocked, suspected := IsBot(60, cfg)
	if blocked {
		t.Error("score 60 should not be blocked")
	}
	if !suspected {
		t.Error("score 60 should be suspected")
	}
}

func TestIsBot_Blocked(t *testing.T) {
	cfg := &config.BotConfiguration{MarkThreshold: 50, BlockThreshold: 80}
	blocked, suspected := IsBot(90, cfg)
	if !blocked {
		t.Error("score 90 should be blocked")
	}
	if !suspected {
		t.Error("score 90 should be suspected")
	}
}
