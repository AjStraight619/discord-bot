package sportsutils

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Teams struct {
	Teams []Team `json:"teams"`
}

type Team struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func LoadTeams() *Teams {
	rootDir, err := filepath.Abs("")
	log.Printf("Root dir: %s", rootDir)
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

func (teams *Teams) FindTeam(query string) (*Team, error) {
	for _, team := range teams.Teams {
		if strings.EqualFold(team.Name, query) || strings.EqualFold(team.ID, query) {
			return &team, nil
		}
	}
	return nil, errors.New("team not found")
}
