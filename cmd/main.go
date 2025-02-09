package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	deps "github.com/AjStraight619/discord-bot/deps"
	"github.com/AjStraight619/discord-bot/internal/bot"
	"github.com/AjStraight619/discord-bot/internal/config"
	"github.com/AjStraight619/discord-bot/internal/messaging"

	"github.com/bwmarrin/discordgo"
)

func main() {
	rootDir := deps.GetProjectRoot()
	binDir := filepath.Join(rootDir, "bin")
	os.MkdirAll(binDir, os.ModePerm)

	deps.EnsureYTDLP(binDir)
	deps.EnsureFFmpeg(binDir)

	config.LoadConfig()

	dg, err := discordgo.New("Bot " + config.AppConfig.DiscordKey)

	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	dg.StateEnabled = true
	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers

	botController := &bot.BotController{
		Session:         dg,
		TimeoutDuration: time.Duration(20) * time.Minute,
	}

	botController.InitCommands()

	dg.AddHandler(botController.MessageHandler)

	// Open a connection to Discord
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error connecting to Discord: %v", err)
	}

	time.Sleep(3 * time.Second)

	fmt.Println("Bot is now running! Press CTRL+C to exit.")

	cm := messaging.InitCron(dg, 15)

	cm.StartCydCron()

	// guild := utils.FindGuildByName(dg, "King's Landing")

	// if guild == nil {
	// 	log.Println("Couldnt find guild")
	// 	return
	// }

	// allMembers, err := members.FetchAllGuildMembers(dg, guild.ID)

	// if err != nil {
	// 	log.Printf("Not able to fetch members: %v", err)
	// }

	// for _, member := range allMembers {
	// 	log.Printf("Member: %s (ID: %s)", member.User.Username, member.User.ID)
	// 	// if strings.EqualFold(member.User.Username, "crispylols") {
	// 	// 	messaging.SendDM(dg, member.User.ID, "Matt did u get this?")
	// 	// }
	// }

	// Graceful shutdown handling
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	// Cleanup
	dg.Close()
}
