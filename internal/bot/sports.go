package bot

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/go-resty/resty/v2"
)

type SportsClient struct {
	APIKey      string
	PlayerStats *PlayerStats
	TeamStats   *TeamStats
}

type SportsQuery struct {
	SportsClient *SportsClient
	APIKey       string
	Team         string
	Name         string
	Season       string
}

type SportsCommand struct{}

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

type PlayerStats struct {
	StatID            int    `json:"StatID"`
	TeamID            int    `json:"TeamID"`
	PlayerID          int    `json:"PlayerID"`
	SeasonType        int    `json:"SeasonType"`
	Season            int    `json:"Season"`
	Name              string `json:"Name"`
	Team              string `json:"Team"`
	Position          string `json:"Position"`
	Started           int    `json:"Started"`
	FanDuelSalary     int    `json:"FanDuelSalary"`
	DraftKingsSalary  int    `json:"DraftKingsSalary"`
	FantasyDataSalary int    `json:"FantasyDataSalary"`
	YahooSalary       int    `json:"YahooSalary"`
	InjuryStatus      string `json:"InjuryStatus"`
	InjuryBodyPart    string `json:"InjuryBodyPart"`
	// InjuryStartDate may be null in the JSON, so we use a pointer.
	InjuryStartDate      *string `json:"InjuryStartDate"`
	InjuryNotes          string  `json:"InjuryNotes"`
	FanDuelPosition      string  `json:"FanDuelPosition"`
	DraftKingsPosition   string  `json:"DraftKingsPosition"`
	YahooPosition        string  `json:"YahooPosition"`
	OpponentRank         int     `json:"OpponentRank"`
	OpponentPositionRank int     `json:"OpponentPositionRank"`
	GlobalTeamID         int     `json:"GlobalTeamID"`
	// FantasyDraftSalary may be null so we use a pointer.
	FantasyDraftSalary   *int   `json:"FantasyDraftSalary"`
	FantasyDraftPosition string `json:"FantasyDraftPosition"`
	GameID               int    `json:"GameID"`
	OpponentID           int    `json:"OpponentID"`
	Opponent             string `json:"Opponent"`
	// You can change these to time.Time if you prefer parsing the date/time strings.
	Day                           string  `json:"Day"`
	DateTime                      string  `json:"DateTime"`
	HomeOrAway                    string  `json:"HomeOrAway"`
	IsGameOver                    bool    `json:"IsGameOver"`
	GlobalGameID                  int     `json:"GlobalGameID"`
	GlobalOpponentID              int     `json:"GlobalOpponentID"`
	Updated                       string  `json:"Updated"`
	Games                         int     `json:"Games"`
	FantasyPoints                 float64 `json:"FantasyPoints"`
	Minutes                       int     `json:"Minutes"`
	Seconds                       int     `json:"Seconds"`
	FieldGoalsMade                float64 `json:"FieldGoalsMade"`
	FieldGoalsAttempted           float64 `json:"FieldGoalsAttempted"`
	FieldGoalsPercentage          float64 `json:"FieldGoalsPercentage"`
	EffectiveFieldGoalsPercentage float64 `json:"EffectiveFieldGoalsPercentage"`
	TwoPointersMade               float64 `json:"TwoPointersMade"`
	TwoPointersAttempted          float64 `json:"TwoPointersAttempted"`
	TwoPointersPercentage         float64 `json:"TwoPointersPercentage"`
	ThreePointersMade             float64 `json:"ThreePointersMade"`
	ThreePointersAttempted        float64 `json:"ThreePointersAttempted"`
	ThreePointersPercentage       float64 `json:"ThreePointersPercentage"`
	FreeThrowsMade                float64 `json:"FreeThrowsMade"`
	FreeThrowsAttempted           float64 `json:"FreeThrowsAttempted"`
	FreeThrowsPercentage          float64 `json:"FreeThrowsPercentage"`
	OffensiveRebounds             float64 `json:"OffensiveRebounds"`
	DefensiveRebounds             float64 `json:"DefensiveRebounds"`
	Rebounds                      float64 `json:"Rebounds"`
	OffensiveReboundsPercentage   float64 `json:"OffensiveReboundsPercentage"`
	DefensiveReboundsPercentage   float64 `json:"DefensiveReboundsPercentage"`
	TotalReboundsPercentage       float64 `json:"TotalReboundsPercentage"`
	Assists                       float64 `json:"Assists"`
	Steals                        float64 `json:"Steals"`
	BlockedShots                  float64 `json:"BlockedShots"`
	Turnovers                     float64 `json:"Turnovers"`
	PersonalFouls                 float64 `json:"PersonalFouls"`
	Points                        float64 `json:"Points"`
	TrueShootingAttempts          float64 `json:"TrueShootingAttempts"`
	TrueShootingPercentage        float64 `json:"TrueShootingPercentage"`
	PlayerEfficiencyRating        float64 `json:"PlayerEfficiencyRating"`
	AssistsPercentage             float64 `json:"AssistsPercentage"`
	StealsPercentage              float64 `json:"StealsPercentage"`
	BlocksPercentage              float64 `json:"BlocksPercentage"`
	TurnOversPercentage           float64 `json:"TurnOversPercentage"`
	UsageRatePercentage           float64 `json:"UsageRatePercentage"`
	FantasyPointsFanDuel          float64 `json:"FantasyPointsFanDuel"`
	FantasyPointsDraftKings       float64 `json:"FantasyPointsDraftKings"`
	FantasyPointsYahoo            float64 `json:"FantasyPointsYahoo"`
	PlusMinus                     float64 `json:"PlusMinus"`
	DoubleDoubles                 float64 `json:"DoubleDoubles"`
	TripleDoubles                 float64 `json:"TripleDoubles"`
	FantasyPointsFantasyDraft     float64 `json:"FantasyPointsFantasyDraft"`
	IsClosed                      bool    `json:"IsClosed"`
	LineupConfirmed               bool    `json:"LineupConfirmed"`
	LineupStatus                  string  `json:"LineupStatus"`
}

type TeamStats struct {
	Team   string `json:"Team"`
	Name   string `json:"Name"`
	Wins   int    `json:"Wins"`
	Losses int    `json:"Losses"`
	Season int    `json:"Season"`
}

func (sc SportsCommand) Execute(b *BotController, msg *discordgo.MessageCreate, options []string) {

	if len(options) < 1 {
		b.displayCmdError(msg.ChannelID, "⚠ Usage: `!sports <team: (e.g: LAL)> <player: (e.g: LeBron James)>`")
		return
	}

	team := options[0]
	name := options[1]

	sportsQuery := b.SportsClient.NewSportsQuery(team, name)

	log.Printf("Sports query: %+v", sportsQuery)

}

func (sc SportsCommand) Help() {

}

func (s SportsClient) NewSportsQuery(team, name string) *SportsQuery {
	return &SportsQuery{
		APIKey: s.APIKey,
		Team:   team,
		Name:   name,
		Season: "2024POST", // Default season
	}
}

func NewSportsClient(apiKey string) *SportsClient {
	return &SportsClient{APIKey: apiKey}
}

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

func (s SportsClient) GetPlayerSeasonStats(season string) ([]PlayerStats, error) {
	client := resty.New()

	url := fmt.Sprintf("https://api.sportsdata.io/v3/nba/stats/json/PlayerSeasonStats/%s?key=%s", season, s.APIKey)

	// Pass the API key in the header as required.
	resp, err := client.R().
		SetHeader("Ocp-Apim-Subscription-Key", s.APIKey).
		Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	// Unmarshal the response body into a slice of PlayerSeason.
	var stats []PlayerStats
	if err := json.Unmarshal(resp.Body(), &stats); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	return stats, nil
}

func FindPlayerSeasonStats(stats []PlayerStats, name string) *PlayerStats {
	for _, p := range stats {
		if p.Name == name {
			return &p
		}
	}
	return nil
}

func GetTotalGamesByTeam(teamName string) (int, error) {
	return 0, nil
}

func (s SportsClient) GetTeamStats(teamName string) *TeamStats {

	url := fmt.Sprintf("https://api.sportsdata.io/v3/nba/scores/json/TeamSeasonStats/2024?key=%s", s.APIKey)

	log.Printf("url: %s", url)

	return &TeamStats{}
}
