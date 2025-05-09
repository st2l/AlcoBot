package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/st2l/AlcoBot/internal/bot"
	"github.com/st2l/AlcoBot/internal/config"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize configuration
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}
	config.SetConfig(cfg)

	// Create a context that will be canceled on SIGINT or SIGTERM
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Initialize and start the bot
	b, err := bot.New()
	if err != nil {
		log.Fatalf("Failed to initialize bot: %v", err)
	}

	b.Start()

	// Keep the bot running until context is canceled
	<-ctx.Done()
	log.Println("Shutting down bot...")
}
