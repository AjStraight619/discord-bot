package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/AjStraight619/discord-bot/internal/apiclients"
	"github.com/bwmarrin/discordgo"
)

// Teams and Team structures remain the same.
type Teams struct {
	Teams []Team `json:"teams"`
}

type Team struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// SportsCommand defines the sports command.
type SportsCommand struct{}

func (sc SportsCommand) Execute(b *BotController, msg *discordgo.MessageCreate, options []string) {
	// Validate that at least a team and a player name are provided.
	if len(options) < 2 {
		b.displayCmdError(msg.ChannelID, "⚠ Usage: `!sports <team> <player> [season]`")
		return
	}

	// Load teams data from file.
	teams := LoadTeams()
	if teams == nil {
		b.displayCmdError(msg.ChannelID, "⚠ Error loading team data.")
		return
	}

	// Create a SportsQuery using the options and loaded teams.
	query, err := NewSportsQuery(options, teams)
	if err != nil {
		b.displayCmdError(msg.ChannelID, fmt.Sprintf("⚠ %s", err.Error()))
		return
	}

	// Determine season and mode.
	// If no season is provided, you can default to a value (for example, "2024").
	season := query.Season
	if season == "" {
		season = "2024" // Adjust as needed.
	}
	// Assume mode "REG" for regular season. (Alternatively, you could support an optional mode parameter.)
	mode := "REG"

	// Call the external API function to get team statistics.
	stats, err := apiclients.GetTeamStatistics(query.TeamID, season, mode)
	if err != nil {
		b.displayCmdError(msg.ChannelID, "Error fetching team statistics.")
		return
	}

	// Send the result to the Discord channel.
	b.Session.ChannelMessageSend(msg.ChannelID, stats)
}

func (sc SportsCommand) Help() string {
	return "!sports <team> <player> [season] - Query for NBA stats. Example: !sports LAL \"LeBron James\" 2024POST"
}

// LoadTeams loads the teams from a JSON file.
func LoadTeams() *Teams {
	// Adjust the path as needed. Here, we assume the data file is at the project root in a 'data' folder.
	rootDir, err := filepath.Abs("../../")
	if err != nil {
		log.Printf("Error getting rootDir: %v", err)
		return nil
	}

	dataPath := filepath.Join(rootDir, "data", "nba_teams.json")
	data, err := os.ReadFile(dataPath)
	if err != nil {
		log.Printf("Error reading file: %v", err)
		return nil
	}

	var teams Teams
	err = json.Unmarshal(data, &teams)
	if err != nil {
		log.Printf("Error unmarshaling data: %v", err)
		return nil
	}

	log.Printf("Teams loaded: %+v", teams)
	return &teams
}

// FindTeam searches for a team by name or ID.
func (teams *Teams) FindTeam(query string) (*Team, error) {
	for _, team := range teams.Teams {
		if strings.EqualFold(team.Name, query) || strings.EqualFold(team.ID, query) {
			return &team, nil
		}
	}
	return nil, errors.New("team not found")
}

// SportsQuery encapsulates a sports query.
type SportsQuery struct {
	TeamID     string
	TeamName   string
	PlayerName string
	Season     string
}

func NewSportsQuery(options []string, teams *Teams) (*SportsQuery, error) {
	teamQuery := options[0]
	playerQuery := options[1]

	team, err := teams.FindTeam(teamQuery)
	if err != nil {
		return nil, err
	}

	query := &SportsQuery{
		TeamID:     team.ID,
		TeamName:   team.Name,
		PlayerName: playerQuery,
	}
	if len(options) >= 3 {
		query.Season = options[2]
	}
	return query, nil
}
