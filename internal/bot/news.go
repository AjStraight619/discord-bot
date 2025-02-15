// package bot

// import (
// 	"fmt"
// 	"log"

// 	"github.com/bwmarrin/discordgo"
// 	"github.com/go-resty/resty/v2"
// )

// // NewsClient struct to handle NewsAPI requests

// // NewsResponse struct to parse API JSON response
// type NewsResponse struct {
// 	Status       string `json:"status"`
// 	TotalResults int    `json:"totalResults"`
// 	Articles     []struct {
// 		Title       string `json:"title"`
// 		Description string `json:"description"`
// 		URL         string `json:"url"`
// 		Source      struct {
// 			Name string `json:"name"`
// 		} `json:"source"`
// 	} `json:"articles"`
// }

// type NewsCommand struct{}

// // Execute fetches and sends the top news headlines for the provided country code.
// func (n NewsCommand) Execute(b *BotController, msg *discordgo.MessageCreate, options []string) {
// 	if len(options) == 0 {
// 		b.Session.ChannelMessageSend(msg.ChannelID, "Please specify a country code. Example: `!news us`")
// 		return
// 	}

// 	country := options[0]
// 	newsChan := make(chan string)

// 	// Fetch news in a separate goroutine.
// 	go b.NewsClient.FetchTopNews(country, newsChan)

// 	// Wait for the response.
// 	newsMessage := <-newsChan

// 	// Send the news message to the Discord channel.
// 	b.Session.ChannelMessageSend(msg.ChannelID, newsMessage)
// }

// // Help returns a usage string for the news command.
// func (n NewsCommand) Help() string {
// 	return "!news <country code> - Displays the top headlines for the specified country."
// }

// // FetchTopNews runs in a goroutine and sends the result to the channel
// func (n *NewsClient) FetchTopNews(country string, newsChan chan<- string) {
// 	url := fmt.Sprintf("https://newsapi.org/v2/top-headlines?country=%s&apiKey=%s", country, n.APIKey)

// 	client := resty.New()
// 	resp, err := client.R().SetResult(&NewsResponse{}).Get(url)

// 	if err != nil {
// 		log.Printf("Error fetching news: %v", err)
// 		newsChan <- "Error fetching news. Please try again."
// 		return
// 	}

// 	news := resp.Result().(*NewsResponse)

// 	if len(news.Articles) == 0 {
// 		newsChan <- "No news articles found for this country."
// 		return
// 	}

// 	// Format news articles
// 	var newsMessage string
// 	for i, article := range news.Articles {
// 		if i >= 5 {
// 			break
// 		}
// 		newsMessage += fmt.Sprintf("**%s** - [%s](%s)\n", article.Title, article.Source.Name, article.URL)
// 	}

// 	log.Printf("fetched news: %s", newsMessage)

// 	newsChan <- newsMessage
// }

// // DisplayNewsResponse is called when `!news` is used
// func (b *BotController) DisplayNewsResponse(options []string, msg *discordgo.MessageCreate) {
// 	if len(options) == 0 {
// 		b.Session.ChannelMessageSend(msg.ChannelID, "Please specify a country code. Example: `!news us`")
// 		return
// 	}

// 	country := options[0]

// 	newsChan := make(chan string)

// 	go b.NewsClient.FetchTopNews(country, newsChan)

// 	newsMessage := <-newsChan

// 	// Send the news message to Discord
// 	b.Session.ChannelMessageSend(msg.ChannelID, newsMessage)
// }

package bot

import (
	"log"

	"github.com/AjStraight619/discord-bot/internal/apiclients"
	"github.com/bwmarrin/discordgo"
)

type NewsCommand struct{}

func (n NewsCommand) Execute(b *BotController, msg *discordgo.MessageCreate, options []string) {
	if len(options) == 0 {
		b.Session.ChannelMessageSend(msg.ChannelID, "Please specify a country code. Example: `!news us`")
		return
	}

	country := options[0]

	// Call the external API function directly.
	newsMessage, err := apiclients.GetTopNews(country)
	if err != nil {
		log.Printf("Error getting news: %v", err)
		b.Session.ChannelMessageSend(msg.ChannelID, "Error fetching news. Please try again.")
		return
	}

	b.Session.ChannelMessageSend(msg.ChannelID, newsMessage)
}

func (n NewsCommand) Help() string {
	return "!news <country code> - Displays the top headlines for the specified country."
}
