package bot

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// BotController holds the state for your bot.
type BotController struct {
	VoiceTextChannelID string
	Session            *discordgo.Session
	NewsClient         *NewsClient
	AIClient           *AIClient
	SportsClient       *SportsClient
	VoiceConn          *discordgo.VoiceConnection
	VoiceGuildMap      map[string]*discordgo.VoiceConnection
	Voices             map[string]*PersonVoice
	isBotInChannel     bool
	musicQueue         []*Song
	isPlaying          bool
	channelName        string
	LastHeard          time.Time
	CommandRegistry    *CommandRegistry
	inactivityTimer    *time.Timer
	TimeoutDuration    time.Duration
	VoiceHandler       *VoiceCommandHandler
}

// MessageHandler handles incoming text commands.
func (b *BotController) MessageHandler(s *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.ID == s.State.User.ID {
		return
	}

	// Reset timer

	// Set the text channel ID for voice responses.
	b.VoiceTextChannelID = msg.ChannelID

	msgContent := strings.TrimSpace(msg.Content)

	if !strings.HasPrefix(msgContent, "!") {
		return // Ignore messages that are not commands
	}

	msgParts := strings.Fields(msgContent)

	var actions []string
	var options []string

	for _, part := range msgParts {
		if strings.HasPrefix(part, "!") {
			actions = append(actions, part)
		} else {
			options = append(options, part)
		}
	}

	log.Printf("Parsed Actions: %v", actions)
	log.Printf("Parsed Options: %v", options)

	// Process each action using the command registry.
	for _, action := range actions {
		if cmd, ok := b.CommandRegistry.Get(action); ok {
			go cmd.Execute(b, msg, options)
		} else {
			log.Printf("Unknown command: %s", action)
			b.displayCmdError(msg.ChannelID, fmt.Sprintf("Unknown command: %s", action))
		}
	}
}

// InitCommands initializes the command registry and registers commands.
func (b *BotController) InitCommands() {
	b.CommandRegistry = NewCommandRegistry()
	// b.CommandRegistry.Register("!ping", PingCommand{})
	b.CommandRegistry.Register("!news", NewsCommand{})
	b.CommandRegistry.Register("!ai", AICommand{})
	b.CommandRegistry.Register("!play", SongCommand{})
	b.CommandRegistry.Register("!listen", ListenCommand{})
	b.CommandRegistry.Register("!join", JoinCommand{})
}

func (m *BotController) displayCmdError(channelID string, msg string) {
	m.Session.ChannelMessageSend(channelID, msg)
}

// joinUserChannel is a helper to join a voice channel.
func (b *BotController) joinUserChannel(guildID, userID string, mute, deafened bool) (*discordgo.VoiceConnection, error) {
	guild, err := b.Session.State.Guild(guildID)
	if err != nil {
		return nil, fmt.Errorf("cannot get guild from state: %w", err)
	}
	for _, vs := range guild.VoiceStates {
		if vs.UserID == userID {
			vc, err := b.Session.ChannelVoiceJoin(guildID, vs.ChannelID, mute, deafened)
			if err != nil {
				return nil, fmt.Errorf("failed to join voice channel: %w", err)
			}
			b.VoiceConn = vc
			// NOTE: Probably better to just attach vc to struct instead of returning it.
			return vc, nil
		}
	}
	return nil, fmt.Errorf("user not in a voice channel")
}

func (b *BotController) ResetTimeout() {
	log.Println("ResetTimeout: Called to restart the inactivity timer.")

	// If there's an existing timer, stop it.
	if b.inactivityTimer != nil {
		log.Println("ResetTimeout: An existing timer was found; attempting to stop it.")
		// Stop returns false if the timer has already expired.
		if !b.inactivityTimer.Stop() {
			log.Println("ResetTimeout: Timer had already expired; draining the timer's channel.")
			// Drain the timer's channel if needed.
			select {
			case <-b.inactivityTimer.C:
				log.Println("ResetTimeout: Drained one value from the timer's channel.")
			default:
				log.Println("ResetTimeout: No value to drain from the timer's channel.")
			}
		} else {
			log.Println("ResetTimeout: Timer stopped successfully.")
		}
	} else {
		log.Println("ResetTimeout: No existing timer found. Creating a new one.")
	}

	// Start a new timer with the configured duration.
	b.inactivityTimer = time.AfterFunc(b.TimeoutDuration, func() {
		log.Println("Timeout reached. Executing timeout action.")
		b.OnTimeout()
	})

	log.Printf("ResetTimeout: New timer started with a timeout duration of %v.\n", b.TimeoutDuration)
}

func (b *BotController) OnTimeout() {
	log.Println("Timeout reached (inactivity)")
	b.LeaveVoiceChannel()
}

func (b *BotController) LeaveVoiceChannel() {
	if b.VoiceConn != nil {
		log.Println("👋 Leaving voice channel...")
		close(b.VoiceConn.OpusRecv)
		b.VoiceConn.Disconnect()
		b.VoiceConn = nil
		b.isBotInChannel = false
		b.Session.ChannelMessageSend(b.VoiceTextChannelID, "✅ Left the voice channel due to inactivity.")
	}
}

func extractCommands(message string) []string {
	re := regexp.MustCompile(`!([a-zA-Z]+)`)
	matches := re.FindAllString(message, -1)
	return matches
}

// MonitorVoiceIdle, LeaveVoiceChannel, and other helper functions can also reside here.
// func (b *BotController) MonitorVoiceIdle() {
// 	ticker := time.NewTicker(10 * time.Second)
// 	defer ticker.Stop()
// 	for range ticker.C {
// 		if b.isBotInChannel && time.Since(b.LastHeard) > 60*time.Second {
// 			log.Println("No audio detected for over 60 seconds. Leaving voice channel...")
// 			b.LeaveVoiceChannel()
// 			return
// 		}
// 	}
// }
