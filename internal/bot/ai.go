package bot

import (
	"fmt"
	"log"
	"strings"

	"github.com/AjStraight619/discord-bot/internal/apiclients"
	"github.com/bwmarrin/discordgo"
)

type AICommand struct{}

func (ai AICommand) Execute(b *BotController, msg *discordgo.MessageCreate, options []string) {
	if len(options) == 0 {
		b.Session.ChannelMessageSend(msg.ChannelID, "Please enter a question! Example: `!ai How does the quadratic formula work?`")
		return
	}

	b.ChatGPTResponse(options, msg)
}

func (ai AICommand) Help() string {
	return "!ai <question> - Ask the AI a question."
}

func (b *BotController) ChatGPTResponse(options []string, msg *discordgo.MessageCreate) {

	query := strings.Join(options, " ")
	response, err := apiclients.GetAIResponse(query)
	if err != nil {
		log.Printf("Error getting AI response: %v", err)
		b.Session.ChannelMessageSend(msg.ChannelID, "Error fetching AI response. Please try again later.")
		return
	}
	b.Session.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("🤖 **ChatGPT:** %s", response))
}
