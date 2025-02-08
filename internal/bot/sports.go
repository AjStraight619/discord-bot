package bot

import (
	"encoding/json"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/go-resty/resty/v2"
)

type SportsClient struct {
	APIKey       string
	PlayerData   *Player
	GameStatData *GameStat
}

type SportsCommand struct{}

type SportsData struct {
}

type Player struct {
	PlayerID  int    `json:"PlayerID"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	Team      string `json:"Team"`
}

type PlayerSeason struct {
	PlayerID  int     `json:"PlayerID"`
	FirstName string  `json:"FirstName"`
	LastName  string  `json:"LastName"`
	Team      string  `json:"Team"`
	Points    float64 `json:"Points"`
	Assists   float64 `json:"Assists"`
	Rebounds  float64 `json:"Rebounds"`
}

type GameStat struct {
	GameID   int     `json:"GameID"`
	Date     string  `json:"Date"`
	Points   float64 `json:"Points"`
	Assists  float64 `json:"Assists"`
	Rebounds float64 `json:"Rebounds"`
}

func (sc SportsCommand) Execute(b *BotController, msg *discordgo.MessageCreate, options []string) {

	if len(options) < 1 {
		b.displayCmdError(msg.ChannelID, "⚠ Usage: `!play <music_link>`")
		return
	}

}

func (sc SportsCommand) Help() {

}

func NewSportsClient(apiKey string) *SportsClient {
	return &SportsClient{APIKey: apiKey}
}

// func (s SportsClient) GetSportData() {
// 	client := resty.New()
// 	url := "https://api.sportsdata.io/v3/nba/scores/json/GamesByDate/2025-FEB-08"
//
// 	resp, err := client.R().
// 		SetHeader("Ocp-Apim-Subscription-Key", s.APIKey). // SportsDataIO requires this header for authentication
// 		Get(url)
//
// 	if err != nil {
// 		log.Fatalf("Error fetching data: %v", err)
// 	}
//
// 	fmt.Println("Response Status Code:", resp.StatusCode())
// 	fmt.Println("Response Body:", resp.String())
//
// }

func (s SportsClient) GetSportData(apiURL string) (*resty.Response, error) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Ocp-Apim-Subscription-Key", s.APIKey).
		Get(apiURL)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s SportsClient) GetPlayerSeasonStats(season string) ([]PlayerSeason, error) {
	client := resty.New()

	// Build the URL using the season parameter. Note that we append the API key as a query parameter.
	url := fmt.Sprintf("https://api.sportsdata.io/v3/nba/stats/json/PlayerSeasonStats/%s?key=%s", season, s.APIKey)

	resp, err := client.R().Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	// Unmarshal the response body into a slice of PlayerSeason.
	var stats []PlayerSeason
	if err := json.Unmarshal(resp.Body(), &stats); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	return stats, nil
}

func FindPlayerSeasonStats(stats []PlayerSeason, firstName, lastName string) *PlayerSeason {
	for _, p := range stats {
		if p.FirstName == firstName && p.LastName == lastName {
			return &p
		}
	}
	return nil
}
