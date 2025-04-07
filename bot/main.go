package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/amiosamu/interview-match-bot/internal/bot"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
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

	// Get database connection string
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Println("DATABASE_URL not set. Quiz functionality will not work properly.")
		// You could set a default for development here
		dbURL = "postgres://postgres:postgres@localhost:5432/interview_bot?sslmode=disable"
	}

	// Connect to the database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Verify the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	log.Println("Connected to database successfully")

	// Create a new bot instance
	interviewBot, err := bot.NewBot(token, db)
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}

	// Enable debug mode for development
	interviewBot.Debug(true)

	// Start the bot
	log.Println("Starting Interview Match Bot...")
	interviewBot.Start()
}
