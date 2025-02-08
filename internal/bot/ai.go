package bot

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sashabaranov/go-openai"
)

type AIClient struct {
	client *openai.Client
}

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

func NewAIClient(apiKey string) *AIClient {
	return &AIClient{
		client: openai.NewClient(apiKey),
	}
}

func (a *AIClient) GetAIResponse(prompt string, responseChan chan<- string) {
	var resp openai.ChatCompletionResponse
	var err error
	maxRetries := 3
	waitTime := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		resp, err = a.client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{Role: "user", Content: prompt},
				},
			},
		)

		if err == nil {
			break
		}

		// Log the error and retry after waiting
		log.Printf("ChatGPT API error (attempt %d): %v", i+1, err)
		time.Sleep(waitTime)
		waitTime *= 2 // Exponential backoff (2s → 4s → 8s)
	}

	// If after retries there's still an error, return a user-friendly message
	if err != nil {
		responseChan <- "⚠ OpenAI API rate limit exceeded. Please try again later."
		return
	}

	responseChan <- resp.Choices[0].Message.Content
}

// ChatGPTResponse handles `!ai` messages
func (b *BotController) ChatGPTResponse(options []string, msg *discordgo.MessageCreate) {
	if len(options) == 0 {
		b.Session.ChannelMessageSend(msg.ChannelID, "Please enter a question! Example: `!ai What is Go?`")
		return
	}

	query := strings.Join(options, " ")

	responseChan := make(chan string)

	// Fetch AI response in a goroutine
	go b.AIClient.GetAIResponse(query, responseChan)

	// Wait for response and send to Discord
	aiResponse := <-responseChan
	b.Session.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("🤖 **ChatGPT:** %s", aiResponse))
}

func (b *BotController) NewsSummaryCommand(options []string, msg *discordgo.MessageCreate) {
	country := options[0]

	newsChan := make(chan string)
	summaryChan := make(chan string)

	go b.NewsClient.FetchTopNews(country, newsChan)

	news := <-newsChan

	go b.AIClient.GetAIResponse(news, summaryChan)

	summary := <-summaryChan

	finalMessage := "**📰 News Summary:**\n" + summary
	b.Session.ChannelMessageSend(msg.ChannelID, finalMessage)
}
