package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tracking/analysis/internal/cache"
	"github.com/tracking/analysis/internal/config"
	"github.com/tracking/analysis/internal/database"
	"github.com/tracking/analysis/internal/handler"
	"github.com/tracking/analysis/internal/repo"
	"github.com/tracking/analysis/internal/rpc"
	"github.com/tracking/analysis/internal/security"
)

func main() {
	configFilename := flag.String("c", "config", "config file name (without extension)")
	configDirs := flag.String("cPath", "./,./configs/", "comma-separated config search paths")
	flag.Parse()

	// Load config
	var cfg config.Config
	if err := config.InitConfiguration(*configFilename, strings.Split(*configDirs, ","), &cfg); err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Set log level
	var level slog.Level
	switch cfg.ServiceConfiguration.LogLevel {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})))

	// Auto-generate RSA keys if missing
	if err := security.EnsureKeyPair(cfg.SecurityConfiguration.RSAPrivateKeyPath, cfg.SecurityConfiguration.RSAPublicKeyPath); err != nil {
		slog.Error("failed to ensure RSA key pair", "error", err)
		os.Exit(1)
	}
	slog.Info("RSA key pair ready")

	// Init database
	db, err := database.Init(cfg.PostgresConfiguration.DSN())
	if err != nil {
		slog.Error("failed to init database", "error", err)
		os.Exit(1)
	}
	slog.Info("database connected and migrated")

	// Init Redis
	rdb, err := cache.Init(cfg.RedisConfiguration.Addr, cfg.RedisConfiguration.Password, cfg.RedisConfiguration.Db)
	if err != nil {
		slog.Error("failed to init redis", "error", err)
		os.Exit(1)
	}
	slog.Info("redis connected")

	// Load RSA key pair
	privKey, pubKey, err := security.LoadKeyPair(cfg.SecurityConfiguration.RSAPrivateKeyPath, cfg.SecurityConfiguration.RSAPublicKeyPath)
	if err != nil {
		slog.Error("failed to load RSA key pair", "error", err)
		os.Exit(1)
	}
	slog.Info("RSA keys loaded", "kid", cfg.SecurityConfiguration.KID)

	// Create repos
	trackerRepo := repo.NewTrackerRepo(db)
	campaignRepo := repo.NewCampaignRepo(db)
	channelRepo := repo.NewChannelRepo(db)
	targetRepo := repo.NewTargetRepo(db)
	siteRepo := repo.NewSiteRepo(db)
	clickRepo := repo.NewClickRepo(db)
	eventRepo := repo.NewEventRepo(db)

	// Set up JSON-RPC dispatcher
	dispatcher := rpc.NewDispatcher()

	// Register admin handlers
	adminHandlers := &rpc.AdminHandlers{
		Config:       &cfg,
		TrackerRepo:  trackerRepo,
		CampaignRepo: campaignRepo,
		ChannelRepo:  channelRepo,
		TargetRepo:   targetRepo,
		SiteRepo:     siteRepo,
	}
	adminHandlers.Register(dispatcher)

	// Register track handlers
	trackHandlers := &rpc.TrackHandlers{
		Config:    &cfg,
		Redis:     rdb,
		PrivKey:   privKey,
		ClickRepo: clickRepo,
		EventRepo: eventRepo,
		SiteRepo:  siteRepo,
	}
	dispatcher.Register("track.collectClick", trackHandlers.CollectClick)
	dispatcher.Register("track.collectEvents", trackHandlers.CollectEvents)

	// Set up tracking HTTP handlers
	trackingHandler := &handler.TrackingHandler{
		Config:     &cfg,
		ClickRepo:  clickRepo,
		TargetRepo: targetRepo,
		PubKey:     pubKey,
		PrivKey:    privKey,
	}

	// Set up Gin router
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(requestIDMiddleware())

	// Routes
	r.POST("/rpc", dispatcher.GinHandler())
	r.GET("/t/:token", trackingHandler.HandleJSTrack)
	r.GET("/r/:token", trackingHandler.HandleRedirectTrack)
	r.GET("/sdk/track.js", trackingHandler.HandleSDK)
	r.GET("/public-keys.json", trackingHandler.HandlePublicKeys)

	addr := fmt.Sprintf(":%s", cfg.ServiceConfiguration.Port)
	slog.Info("starting server", "addr", addr)
	if err := r.Run(addr); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}

func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}
