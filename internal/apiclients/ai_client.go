package apiclients

import (
	"context"
	"log"
	"time"

	"github.com/AjStraight619/discord-bot/internal/config"
	"github.com/sashabaranov/go-openai"
)

// GetAIResponse calls the OpenAI API using the global configuration and returns the response.
func GetAIResponse(prompt string) (string, error) {
	client := openai.NewClient(config.AppConfig.OpenAIKey)
	var resp openai.ChatCompletionResponse
	var err error
	maxRetries := 3
	waitTime := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		resp, err = client.CreateChatCompletion(
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
		log.Printf("OpenAI API error (attempt %d): %v", i+1, err)
		time.Sleep(waitTime)
		waitTime *= 2
	}

	if err != nil {
		return "⚠ OpenAI API rate limit exceeded. Please try again later.", err
	}
	return resp.Choices[0].Message.Content, nil
}
