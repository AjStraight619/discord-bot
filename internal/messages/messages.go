package messages

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Messages struct holds the Discord session instance
type Messages struct {
	Session    *discordgo.Session
	NewsClient *NewsClient
	AIClient   *AIClient
}

// MessageHandler processes incoming messages
func (m *Messages) MessageHandler(s *discordgo.Session, msg *discordgo.MessageCreate) {
	// Ignore bot's own messages
	if msg.Author.ID == s.State.User.ID {
		return
	}

	msgContent := strings.TrimSpace(msg.Content)
	msgParts := strings.Fields(msgContent) // Splits input into words

	var actions []string
	var options []string

	// Separate actions and options
	for _, part := range msgParts {
		if strings.HasPrefix(part, "!") {
			actions = append(actions, part)
		} else {
			options = append(options, part)
		}
	}

	log.Printf("Parsed Actions: %v", actions)
	log.Printf("Parsed Options: %v", options)

	// Check if no actions were provided
	if len(actions) == 0 {
		m.Session.ChannelMessageSend(msg.ChannelID, "⚠ No command provided. Use `!help` to see available commands.")
		return
	}

	// Execute commands in order
	for _, action := range actions {
		switch action {
		case "!ping":
			m.Session.ChannelMessageSend(msg.ChannelID, "Pong! 🏓")
		case "!news":
			// ✅ Handle missing country argument
			if len(options) == 0 {
				m.Session.ChannelMessageSend(msg.ChannelID, "⚠ Please specify a country code. Example: `!news us`")
				return
			}
			go m.DisplayNewsResponse(options, msg)

		case "!ai":
			// ✅ Handle missing text argument
			if len(options) == 0 {
				m.Session.ChannelMessageSend(msg.ChannelID, "⚠ Please enter a question. Example: `!ai What is Go?`")
				return
			}
			go m.ChatGPTResponse(options, msg)

		case "!ai!news": // ✅ Handles `!ai !news` as a combined action
			// ✅ Ensure country argument is provided
			if len(options) == 0 {
				m.Session.ChannelMessageSend(msg.ChannelID, "⚠ Please specify a country code for news. Example: `!ai !news us`")
				return
			}
			go m.NewsSummaryCommand(options, msg)

		default:
			log.Printf("Unknown command: %s", action)
			m.Session.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("⚠ Unknown command: %s", action))
		}
	}

}

func (m *Messages) NewsSummaryCommand(options []string, msg *discordgo.MessageCreate) {
	country := options[0]

	newsChan := make(chan string)
	summaryChan := make(chan string)

	go m.NewsClient.FetchTopNews(country, newsChan)

	news := <-newsChan

	go m.AIClient.GetAIResponse(news, summaryChan)

	summary := <-summaryChan

	finalMessage := "**📰 News Summary:**\n" + summary
	m.Session.ChannelMessageSend(msg.ChannelID, finalMessage)
}
