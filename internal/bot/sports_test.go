package bot

import (
	"log"
	"strings"
	"testing"

	"github.com/AjStraight619/discord-bot/internal/apiclients"
	"github.com/AjStraight619/discord-bot/internal/config"
)

func TestGetPlayerStats(t *testing.T) {
	config.LoadConfig()

	teams := LoadTeams()

	if teams == nil {
		t.Fatalf("Failed to load teams")
	}

	lakers, err := teams.FindTeam("Lakers")

	if err != nil {
		t.Fatalf("Failed to find team")
	}

	teamStats, err := apiclients.GetTeamStatistics(lakers.ID, "2024", "REG")

	if err != nil {
		t.Fatalf("Failed to get team stats: %v", err)
	}

	log.Printf("Team Stats: %+v", teamStats)

	// Print specific fields.
	log.Printf("Team: %s", teamStats.Name)
	log.Printf("Market: %s", teamStats.Market)
	log.Printf("Season: %d %s", teamStats.Season.Year, teamStats.Season.Type)
	log.Printf("Games Played: %d", teamStats.OwnRecord.Total.GamesPlayed)

	for _, player := range teamStats.Players {
		if strings.EqualFold("LeBron James", player.FullName) {
			log.Printf("Averages for %s:", player.FullName)
			log.Printf("Minutes: %.2f", player.Averages.Minutes)
			log.Printf("Points: %.2f", player.Averages.Points)
			log.Printf("Offensive Rebounds: %.2f", player.Averages.OffRebounds)
			log.Printf("Defensive Rebounds: %.2f", player.Averages.DefRebounds)
			log.Printf("Rebounds: %.2f", player.Averages.Rebounds)
			log.Printf("Assists: %.2f", player.Averages.Assists)
			log.Printf("Steals: %.2f", player.Averages.Steals)
			log.Printf("Blocks: %.2f", player.Averages.Blocks)
			log.Printf("Turnovers: %.2f", player.Averages.Turnovers)
		}
	}
}
