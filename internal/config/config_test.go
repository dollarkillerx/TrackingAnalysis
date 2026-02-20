package config

import (
	"os"
	"path/filepath"
	"testing"
)

const validTOML = `
[ServiceConfiguration]
Port = "8080"
Debug = false
ExportURL = "https://example.com/export"
LogLevel = "info"

[PostgresConfiguration]
Host = "localhost"
Port = 5432
User = "testuser"
Password = "testpass"
DBName = "tracking"
SSLMode = false
TimeZone = "UTC"

[RedisConfiguration]
Addr = "localhost:6379"
Password = "redispass"
Db = 2

[AdminConfiguration]
Username = "admin"
Password = "secret"

[SecurityConfiguration]
TokenSecret = "my-secret-key"
TSWindowSeconds = 30
NonceTTLSeconds = 60
DedupSeconds = 3600
RSAPrivateKeyPemPath = "/keys/priv.pem"
RSAPublicKeyPemPath = "/keys/pub.pem"
KID = "key-1"

[RateLimitConfiguration]
PerIPPerMinute = 100
PerIPUAPerMinute = 60
PerTrackerIPPerMinute = 30

[BotConfiguration]
MarkThreshold = 50
BlockThreshold = 80
BlockMode = "flag"
`

func TestInitConfiguration_ValidConfig(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "config.toml"), []byte(validTOML), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	var cfg Config
	if err := InitConfiguration("config", []string{dir}, &cfg); err != nil {
		t.Fatalf("InitConfiguration: %v", err)
	}

	if cfg.ServiceConfiguration.Port != "8080" {
		t.Errorf("ServiceConfiguration.Port = %q, want %q", cfg.ServiceConfiguration.Port, "8080")
	}
	if cfg.ServiceConfiguration.ExportURL != "https://example.com/export" {
		t.Errorf("ServiceConfiguration.ExportURL = %q", cfg.ServiceConfiguration.ExportURL)
	}
	if cfg.PostgresConfiguration.Host != "localhost" {
		t.Errorf("PostgresConfiguration.Host = %q, want %q", cfg.PostgresConfiguration.Host, "localhost")
	}
	if cfg.PostgresConfiguration.Port != 5432 {
		t.Errorf("PostgresConfiguration.Port = %d, want 5432", cfg.PostgresConfiguration.Port)
	}
	if cfg.PostgresConfiguration.DBName != "tracking" {
		t.Errorf("PostgresConfiguration.DBName = %q", cfg.PostgresConfiguration.DBName)
	}
	if cfg.RedisConfiguration.Addr != "localhost:6379" {
		t.Errorf("RedisConfiguration.Addr = %q", cfg.RedisConfiguration.Addr)
	}
	if cfg.RedisConfiguration.Db != 2 {
		t.Errorf("RedisConfiguration.Db = %d, want 2", cfg.RedisConfiguration.Db)
	}
	if cfg.SecurityConfiguration.TokenSecret != "my-secret-key" {
		t.Errorf("SecurityConfiguration.TokenSecret = %q", cfg.SecurityConfiguration.TokenSecret)
	}
	if cfg.SecurityConfiguration.TSWindowSeconds != 30 {
		t.Errorf("SecurityConfiguration.TSWindowSeconds = %d, want 30", cfg.SecurityConfiguration.TSWindowSeconds)
	}
	if cfg.RateLimitConfiguration.PerIPPerMinute != 100 {
		t.Errorf("RateLimitConfiguration.PerIPPerMinute = %d, want 100", cfg.RateLimitConfiguration.PerIPPerMinute)
	}
	if cfg.BotConfiguration.MarkThreshold != 50 {
		t.Errorf("BotConfiguration.MarkThreshold = %d, want 50", cfg.BotConfiguration.MarkThreshold)
	}
	if cfg.BotConfiguration.BlockMode != "flag" {
		t.Errorf("BotConfiguration.BlockMode = %q, want %q", cfg.BotConfiguration.BlockMode, "flag")
	}
}

func TestInitConfiguration_MissingFile(t *testing.T) {
	var cfg Config
	err := InitConfiguration("config", []string{"/nonexistent/path"}, &cfg)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestInitConfiguration_InvalidTOML(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "config.toml"), []byte("this is [[[not valid toml"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	var cfg Config
	err := InitConfiguration("config", []string{dir}, &cfg)
	if err == nil {
		t.Fatal("expected error for invalid TOML")
	}
}

func TestInitConfiguration_PartialConfig(t *testing.T) {
	dir := t.TempDir()
	partial := `
[ServiceConfiguration]
Port = "9090"
`
	if err := os.WriteFile(filepath.Join(dir, "config.toml"), []byte(partial), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	var cfg Config
	if err := InitConfiguration("config", []string{dir}, &cfg); err != nil {
		t.Fatalf("InitConfiguration: %v", err)
	}
	if cfg.ServiceConfiguration.Port != "9090" {
		t.Errorf("ServiceConfiguration.Port = %q, want %q", cfg.ServiceConfiguration.Port, "9090")
	}
	// Unset fields should be zero values
	if cfg.PostgresConfiguration.Host != "" {
		t.Errorf("PostgresConfiguration.Host = %q, want empty", cfg.PostgresConfiguration.Host)
	}
	if cfg.RateLimitConfiguration.PerIPPerMinute != 0 {
		t.Errorf("RateLimitConfiguration.PerIPPerMinute = %d, want 0", cfg.RateLimitConfiguration.PerIPPerMinute)
	}
	if cfg.BotConfiguration.MarkThreshold != 0 {
		t.Errorf("BotConfiguration.MarkThreshold = %d, want 0", cfg.BotConfiguration.MarkThreshold)
	}
}

func TestPostgresConfiguration_DSN(t *testing.T) {
	p := PostgresConfiguration{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "secret",
		DBName:   "mydb",
		SSLMode:  false,
		TimeZone: "UTC",
	}
	dsn := p.DSN()
	expected := "host=localhost user=postgres password=secret dbname=mydb port=5432 TimeZone=UTC sslmode=disable"
	if dsn != expected {
		t.Errorf("DSN() = %q, want %q", dsn, expected)
	}
}

func TestPostgresConfiguration_DSN_SSLEnabled(t *testing.T) {
	p := PostgresConfiguration{
		Host:     "db.example.com",
		Port:     5433,
		User:     "admin",
		Password: "pass",
		DBName:   "prod",
		SSLMode:  true,
	}
	dsn := p.DSN()
	expected := "host=db.example.com user=admin password=pass dbname=prod port=5433"
	if dsn != expected {
		t.Errorf("DSN() = %q, want %q", dsn, expected)
	}
}
