package bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// BotController holds the state for your bot.
type BotController struct {
	VoiceTextChannelID string
	Session            *discordgo.Session
	VoiceConn          *discordgo.VoiceConnection
	isBotInChannel     bool
	musicQueue         []*Song
	isPlaying          bool
	LastHeard          time.Time
	CommandRegistry    *CommandRegistry
	inactivityTimer    *time.Timer
	TimeoutDuration    time.Duration
	VoiceHandler       *VoiceCommandHandler
}

func (b *BotController) MessageHandler(s *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.ID == s.State.User.ID {
		return
	}

	// Set the text channel ID for voice responses.
	b.VoiceTextChannelID = msg.ChannelID

	msgContent := strings.TrimSpace(msg.Content)

	if !strings.HasPrefix(msgContent, "!") {
		return // Ignore messages that are not commands
	}

	b.ResetTimeout()

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
	b.CommandRegistry.Register("!news", NewsCommand{})
	b.CommandRegistry.Register("!ai", AICommand{})
	b.CommandRegistry.Register("!play", SongCommand{})
	b.CommandRegistry.Register("!listen", ListenCommand{})
	b.CommandRegistry.Register("!join", JoinCommand{})
	b.CommandRegistry.Register("!leave", LeaveCommand{})
	b.CommandRegistry.Register("!sports", SportsCommand{})
	b.CommandRegistry.Register("!timeout", TimeoutCommand{})
}

func (b *BotController) displayCmdError(channelID string, msg string) {
	b.Session.ChannelMessageSend(channelID, msg)
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
			b.isBotInChannel = true
			return vc, nil
		}
	}
	return nil, fmt.Errorf("user not in a voice channel")
}

// func extractCommands(message string) []string {
// 	re := regexp.MustCompile(`!([a-zA-Z]+)`)
// 	matches := re.FindAllString(message, -1)
// 	return matches
// }

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
