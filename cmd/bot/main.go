package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sub-cos-counter/internal/bot"
	"sub-cos-counter/internal/config"
	"sub-cos-counter/internal/repository"
	"sub-cos-counter/internal/services"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Log configuration info
	log.Printf("Starting %s v%s in %s mode", 
		cfg.App.Name, cfg.App.Version, cfg.App.Environment)

	if cfg.IsDebugMode() {
		log.Printf("Debug mode enabled")
	}

	if cfg.IsPersonalBot() {
		log.Printf("Personal bot mode - allowed user: %d", cfg.Telegram.AllowedUser)
	}

	// Connect to database
	ctx := context.Background()
	
	// Create connection pool configuration
	poolConfig, err := pgxpool.ParseConfig(cfg.GetDatabaseURL())
	if err != nil {
		log.Fatalf("Failed to parse database URL: %v", err)
	}

	// Configure connection pool settings
	poolConfig.MaxConns = int32(cfg.Database.MaxConnections)
	poolConfig.MaxConnIdleTime = time.Duration(cfg.Database.MaxIdleTime) * time.Minute
	poolConfig.MaxConnLifetime = time.Duration(cfg.Database.ConnMaxLifetime) * time.Minute

	// Create database pool
	dbPool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Test database connection
	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to database")

	// Initialize repositories
	subscriptionRepo := repository.NewSubscriptionRepository(dbPool)
	paymentRepo := repository.NewPaymentRepository(dbPool)

	// Initialize services
	subscriptionService := services.NewSubscriptionService(subscriptionRepo, paymentRepo)
	analyticsService := services.NewAnalyticsService(paymentRepo, subscriptionRepo)

	// Initialize bot
	telegramBot, err := bot.NewBot(cfg, subscriptionService, analyticsService)
	if err != nil {
		log.Fatalf("Failed to initialize bot: %v", err)
	}

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-c
		log.Printf("Received signal %v, shutting down gracefully...", sig)
		
		// Give some time for cleanup
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Stop bot
		telegramBot.Stop()
		
		// Close database connections
		dbPool.Close()
		
		log.Println("Shutdown complete")
		
		select {
		case <-shutdownCtx.Done():
			log.Println("Shutdown timeout exceeded")
		default:
		}
		
		os.Exit(0)
	}()

	// Start bot
	log.Printf("Starting %s...", cfg.App.Name)
	if cfg.IsDevelopment() {
		log.Println("Bot is running in development mode")
		log.Printf("Send /start to the bot to begin!")
	}
	
	telegramBot.Start()
}