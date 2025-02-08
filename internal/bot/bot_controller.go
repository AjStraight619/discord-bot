package bot

//
// import (
// 	"fmt"
// 	"log"
// 	"regexp"
// 	"strings"
// 	"time"
//
// 	"github.com/bwmarrin/discordgo"
// )
//
// // Messages struct holds the Discord session instance
// type BotController struct {
// 	VoiceTextChannelID string
// 	Session            *discordgo.Session
// 	NewsClient         *NewsClient
// 	AIClient           *AIClient
// 	VoiceConn          *discordgo.VoiceConnection
// 	VoiceGuildMap      map[string]*discordgo.VoiceConnection
// 	Voices             map[string]*PersonVoice
// 	isBotInChannel     bool
// 	musicQueue         []*Song
// 	isPlaying          bool
// 	channelName        string
// 	LastHeard          time.Time
// }
//
// func (b *BotController) MessageHandler(s *discordgo.Session, msg *discordgo.MessageCreate) {
// 	if msg.Author.ID == s.State.User.ID {
// 		return
// 	}
//
// 	// Set the VoiceTextChannelID to the channel where the command was received.
// 	b.VoiceTextChannelID = msg.ChannelID
//
// 	msgContent := strings.TrimSpace(msg.Content)
// 	msgParts := strings.Fields(msgContent)
//
// 	var actions []string
// 	var options []string
//
// 	// Separate actions and options
// 	for _, part := range msgParts {
// 		if strings.HasPrefix(part, "!") {
// 			actions = append(actions, part)
// 		} else {
// 			options = append(options, part)
// 		}
// 	}
//
// 	log.Printf("Parsed Actions: %v", actions)
// 	log.Printf("Parsed Options: %v", options)
//
// 	// Execute commands in order
// 	for _, action := range actions {
// 		switch action {
// 		case "!ping":
// 			b.displayCmdError(msg.ChannelID, "Pong! 🏓")
// 		case "!news":
// 			if len(options) == 0 {
// 				b.displayCmdError(msg.ChannelID, "⚠ Please specify a country code. Example: `!news us`")
// 				return
// 			}
// 			go b.DisplayNewsResponse(options, msg)
// 		case "!ai":
// 			if len(options) == 0 {
// 				b.displayCmdError(msg.ChannelID, "⚠ Please enter a question. Example: `!ai What is Go?`")
// 				return
// 			}
// 			go b.ChatGPTResponse(options, msg)
// 		case "!ai!news":
// 			if len(options) == 0 {
// 				b.displayCmdError(msg.ChannelID, "⚠ Please specify a country code for news. Example: `!ai!news us`")
// 				return
// 			}
// 			go b.NewsSummaryCommand(options, msg)
//
// 		case "!join":
//
// 			// TODO: Probably should just put vc on struct directly in function
//
// 			vc, err := b.joinUserChannel(msg.GuildID, msg.Author.ID, false)
//
// 			if err != nil {
// 				b.displayCmdError(msg.ChannelID, "⚠ Failed to join voice channel.")
// 				return
// 			}
//
// 			if vc != nil {
// 				b.VoiceConn = vc
// 				b.isBotInChannel = true
// 				go b.VoiceToTextStream()
// 			}
//
// 		case "!listen":
//
// 			vc, err := b.joinUserChannel(msg.GuildID, msg.Author.ID, false)
// 			if err != nil {
// 				log.Printf("Error joining channel for listen: %v", err)
// 				return
// 			}
//
// 			if vc != nil {
// 				b.VoiceConn = vc
// 				b.isBotInChannel = true
// 				// Start only listening to incoming audio.
// 				go b.ListenVoice(msg)
// 			}
//
// 		case "!deafen":
// 			_, err := b.joinUserChannel(msg.GuildID, msg.Author.ID, true)
// 			if err != nil {
// 				log.Printf("Something went wrong in deafen: %v", err)
// 				return
// 			}
//
// 		case "!leave":
// 			if b.VoiceConn != nil {
// 				log.Println("👋 Leaving voice channel...")
//
// 				// Close the OpusRecv channel to signal Listen() to stop
// 				close(b.VoiceConn.OpusRecv)
//
// 				// Disconnect from voice channel
// 				b.VoiceConn.Disconnect()
// 				b.VoiceConn = nil
// 				b.isBotInChannel = false
//
// 				b.displayCmdError(msg.ChannelID, "✅ Left the voice channel.")
// 			} else {
// 				b.displayCmdError(msg.ChannelID, "⚠ I'm not in a voice channel.")
// 			}
//
// 		case "!play":
// 			if len(options) < 1 {
// 				b.displayCmdError(msg.ChannelID, "⚠ Usage: `!play <music_link>`")
// 			}
//
// 			go b.Play(options, msg)
//
// 		default:
// 			log.Printf("Unknown command: %s", action)
// 			return
// 		}
// 	}
//
// }
//
// func (m *BotController) displayCmdError(channelID string, msg string) {
// 	m.Session.ChannelMessageSend(channelID, msg)
// }
//
// func (b *BotController) joinUserChannel(guildID, userID string, deafened bool) (*discordgo.VoiceConnection, error) {
// 	// First, try to get the guild from the session State
// 	guild, err := b.Session.State.Guild(guildID)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot get guild from state: %w", err)
// 	}
//
// 	// Check every VoiceState in the guild to find the one that matches userID
// 	for _, vs := range guild.VoiceStates {
// 		if vs.UserID == userID {
// 			// Found the channel the user is in
//
// 			vc, err := b.Session.ChannelVoiceJoin(guildID, vs.ChannelID, true, deafened)
// 			if err != nil {
// 				return nil, fmt.Errorf("failed to join voice channel: %w", err)
// 			}
// 			return vc, nil
// 		}
// 	}
//
// 	return nil, fmt.Errorf("user not in a voice channel")
// }
//
// func (b *BotController) MonitorVoiceIdle() {
// 	ticker := time.NewTicker(10 * time.Second)
// 	defer ticker.Stop()
//
// 	for range ticker.C {
// 		if b.isBotInChannel && time.Since(b.LastHeard) > 60*time.Second {
// 			log.Println("No audio detected for over 60 seconds. Leaving voice channel...")
// 			b.LeaveVoiceChannel()
// 			return
// 		}
// 	}
// }
//
// func (b *BotController) LeaveVoiceChannel() {
// 	if b.VoiceConn != nil {
// 		log.Println("👋 Leaving voice channel...")
// 		// Close the OpusRecv channel to signal ListenVoice to stop.
// 		close(b.VoiceConn.OpusRecv)
// 		// Disconnect from the voice channel.
// 		b.VoiceConn.Disconnect()
// 		b.VoiceConn = nil
// 		b.isBotInChannel = false
// 		// Optionally, send a message to a Discord text channel:
// 		b.Session.ChannelMessageSend(b.VoiceTextChannelID, "✅ Left the voice channel due to inactivity.")
// 	}
// }
//
// func extractCommands(message string) []string {
// 	re := regexp.MustCompile(`!([a-zA-Z]+)`)
// 	matches := re.FindAllString(message, -1)
// 	return matches
// }
