package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/amiosamu/interview-match-bot/internal/bot"
)

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it. Will use environment variables.")
	}

	// Get the bot token from environment variable
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatalf("TELEGRAM_BOT_TOKEN environment variable not set")
	}

	// Create a new bot instance
	interviewBot, err := bot.NewBot(token)
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}

	// Enable debug mode for development
	interviewBot.Debug(true)

	// Start the bot
	log.Println("Starting Interview Match Bot...")
	interviewBot.Start()
}