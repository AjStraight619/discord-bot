package main

import (
	"discord-bot/internal/bot"
	"discord-bot/utils"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	rootDir := utils.GetProjectRoot()
	binDir := filepath.Join(rootDir, "bin")
	os.MkdirAll(binDir, os.ModePerm)

	utils.EnsureYTDLP(binDir)
	utils.EnsureFFmpeg(binDir)

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

	newsClient := &bot.NewsClient{APIKey: NewsAPI}
	aiClient := bot.NewAIClient(OpenAIAPI)

	msgHandler := &bot.BotController{
		Session:       dg,
		NewsClient:    newsClient,
		AIClient:      aiClient,
		VoiceGuildMap: nil,
	}

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
