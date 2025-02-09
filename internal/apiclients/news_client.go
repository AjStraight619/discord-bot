package apiclients

import (
	"fmt"
	"log"

	"github.com/AjStraight619/discord-bot/internal/config"
	"github.com/go-resty/resty/v2"
)

// NewsResponse is used to parse the JSON response from the News API.
type NewsResponse struct {
	Status       string `json:"status"`
	TotalResults int    `json:"totalResults"`
	Articles     []struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		URL         string `json:"url"`
		Source      struct {
			Name string `json:"name"`
		} `json:"source"`
	} `json:"articles"`
}

// GetTopNews fetches the top news headlines for a given country using the News API.
// It uses the API key from the global configuration.
func GetTopNews(country string) (string, error) {
	// Get the API key from the global configuration.
	apiKey := config.AppConfig.NewsKey
	url := fmt.Sprintf("https://newsapi.org/v2/top-headlines?country=%s&apiKey=%s", country, apiKey)

	client := resty.New()
	resp, err := client.R().SetResult(&NewsResponse{}).Get(url)
	if err != nil {
		log.Printf("Error fetching news: %v", err)
		return "", err
	}

	news := resp.Result().(*NewsResponse)
	if len(news.Articles) == 0 {
		return "No news articles found for this country.", nil
	}

	var newsMessage string
	for i, article := range news.Articles {
		if i >= 5 {
			break
		}
		newsMessage += fmt.Sprintf("**%s** - [%s](%s)\n", article.Title, article.Source.Name, article.URL)
	}

	log.Printf("Fetched news: %s", newsMessage)
	return newsMessage, nil
}
