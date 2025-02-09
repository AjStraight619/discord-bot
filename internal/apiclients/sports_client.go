package apiclients

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/AjStraight619/discord-bot/internal/config"
)

// GetTeamStatistics fetches team statistics for a given team ID, season, and mode.
// It uses the Sports API key from the global configuration.
func GetTeamStatistics(teamID, season, mode string) (string, error) {
	// Construct the URL based on your API's documentation.
	// For example, if mode is "REG" (regular season) and season is "2024":
	url := fmt.Sprintf("https://api.sportradar.com/nba/trial/v8/en/seasons/%s/%s/teams/%s/statistics.json?api_key=%s",
		season, mode, teamID, config.AppConfig.SportsKey)

	res, err := http.Get(url)
	if err != nil {
		log.Printf("Error performing request: %v", err)
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		return "", err
	}

	// Optionally, you can parse the JSON response here. For now, we return the raw JSON.
	return string(body), nil
}
