package main

import (
	"discord-bot/internal/messages"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Token := os.Getenv("DISCORD_KEY")
	NewsAPI := os.Getenv("NEWS_KEY")
	OpenAIAPI := os.Getenv("OPENAI_KEY")

	if Token == "" || NewsAPI == "" || OpenAIAPI == "" {
		log.Fatal("Missing environment variables.")
	}

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	newsClient := &messages.NewsClient{APIKey: NewsAPI}
	aiClient := messages.NewAIClient(OpenAIAPI)

	msgHandler := &messages.Messages{
		Session:    dg,
		NewsClient: newsClient,
		AIClient:   aiClient, // ✅ Pass AIClient only once
	}

	// Register message handler
	dg.AddHandler(msgHandler.MessageHandler)

	// Open a connection to Discord
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error connecting to Discord: %v", err)
	}

	fmt.Println("Bot is now running! Press CTRL+C to exit.")

	// Graceful shutdown handling
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	// Cleanup
	dg.Close()
}
