package apiclients

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/AjStraight619/discord-bot/internal/config"
	"github.com/AjStraight619/discord-bot/internal/models"
)

// GetTeamStatistics fetches team statistics for a given team ID, season, and mode.
// It uses the Sports API key from the global configuration.
func GetTeamStatistics(teamID, season, mode string) (*models.SRTeam, error) {
	// Construct the URL using the API key from your config.
	log.Printf("API KEY: %s", config.AppConfig.SportsKey)
	url := fmt.Sprintf("https://api.sportradar.com/nba/trial/v8/en/seasons/%s/%s/teams/%s/statistics.json?api_key=%s",
		season, mode, teamID, config.AppConfig.SportsKey)

	res, err := http.Get(url)
	if err != nil {
		log.Printf("Error performing request: %v", err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		return nil, err
	}

	// Unmarshal the JSON response into an SRTeam struct.
	var teamStats models.SRTeam
	err = json.Unmarshal(body, &teamStats)
	if err != nil {
		log.Printf("Error unmarshaling JSON: %v", err)
		return nil, err
	}

	return &teamStats, nil
}
