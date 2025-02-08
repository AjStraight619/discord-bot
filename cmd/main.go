package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/AjStraight619/discord-bot/deps"
	"github.com/AjStraight619/discord-bot/internal/bot"

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
	SportsDataIO := os.Getenv("SPORTS_DATA_IO_KEY")

	if Token == "" || NewsAPI == "" || OpenAIAPI == "" {
		log.Fatal("Missing environment variables.")
	}

	dg, err := discordgo.New("Bot " + Token)

	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	newsClient := &bot.NewsClient{APIKey: NewsAPI}
	aiClient := bot.NewAIClient(OpenAIAPI)
	sportsClient := &bot.SportsClient{APIKey: SportsDataIO}

	botController := &bot.BotController{
		Session:         dg,
		NewsClient:      newsClient,
		AIClient:        aiClient,
		SportsClient:    sportsClient,
		VoiceGuildMap:   nil,
		TimeoutDuration: 20 * time.Second,
	}

	// Initialize the command registry
	botController.InitCommands()

	dg.AddHandler(botController.MessageHandler)

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
