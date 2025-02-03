package bot

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/go-resty/resty/v2"
)

// NewsClient struct to handle NewsAPI requests
type NewsClient struct {
	APIKey string
}

// NewsResponse struct to parse API JSON response
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

// FetchTopNews runs in a goroutine and sends the result to the channel
func (n *NewsClient) FetchTopNews(country string, newsChan chan<- string) {
	url := fmt.Sprintf("https://newsapi.org/v2/top-headlines?country=%s&apiKey=%s", country, n.APIKey)

	client := resty.New()
	resp, err := client.R().SetResult(&NewsResponse{}).Get(url)

	if err != nil {
		log.Printf("Error fetching news: %v", err)
		newsChan <- "Error fetching news. Please try again."
		return
	}

	news := resp.Result().(*NewsResponse)

	if len(news.Articles) == 0 {
		newsChan <- "No news articles found for this country."
		return
	}

	// Format news articles
	var newsMessage string
	for i, article := range news.Articles {
		if i >= 5 {
			break
		}
		newsMessage += fmt.Sprintf("**%s** - [%s](%s)\n", article.Title, article.Source.Name, article.URL)
	}

	log.Printf("fetched news: %s", newsMessage)

	newsChan <- newsMessage
}

// DisplayNewsResponse is called when `!news` is used
func (b *BotController) DisplayNewsResponse(options []string, msg *discordgo.MessageCreate) {
	if len(options) == 0 {
		b.Session.ChannelMessageSend(msg.ChannelID, "Please specify a country code. Example: `!news us`")
		return
	}

	country := options[0]

	newsChan := make(chan string)

	go b.NewsClient.FetchTopNews(country, newsChan)

	newsMessage := <-newsChan

	// Send the news message to Discord
	b.Session.ChannelMessageSend(msg.ChannelID, newsMessage)
}
