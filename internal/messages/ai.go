package messages

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
				Model: openai.GPT3Dot5Turbo, // Use GPT-3.5 by default
				Messages: []openai.ChatCompletionMessage{
					{Role: "user", Content: prompt},
				},
			},
		)

		if err == nil {
			break // ✅ Success, exit retry loop
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
func (m *Messages) ChatGPTResponse(options []string, msg *discordgo.MessageCreate) {
	if len(options) == 0 {
		m.Session.ChannelMessageSend(msg.ChannelID, "Please enter a question! Example: `!ai What is Go?`")
		return
	}

	query := strings.Join(options, " ")
	responseChan := make(chan string)

	// Fetch AI response in a goroutine
	go m.AIClient.GetAIResponse(query, responseChan)

	// Wait for response and send to Discord
	aiResponse := <-responseChan
	m.Session.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("🤖 **ChatGPT:** %s", aiResponse))
}
