package bot

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
)

func TestGetSportData(t *testing.T) {
	rootDir, err := filepath.Abs("../../")
	if err != nil {
		log.Fatal("Error finding root directory:", err)
	}

	envPath := filepath.Join(rootDir, ".env")
	err = godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading .env file from %s", envPath)
	}

	SportsDataIO := os.Getenv("SPORTS_DATA_IO_KEY")
	if SportsDataIO == "" {
		t.Fatal("SPORTS_DATA_IO_KEY is not set in .env")
	}

	sportsClient := SportsClient{APIKey: SportsDataIO}
	resp, err := sportsClient.GetSportData("https://api.sportsdata.io/v3/nba/scores/json/GamesByDate/2025-FEB-08")
	if err != nil {
		t.Fatalf("Error fetching sports data: %v", err)
	}

	// Check HTTP response code
	if resp.StatusCode() != 200 {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode())
	}

	// Print actual response for debugging
	t.Logf("Response Body: %s", resp.String())
}

func TestGetLeBronSeasonStats(t *testing.T) {
	// Load environment variables from the .env file.
	rootDir, err := filepath.Abs("../../")
	if err != nil {
		t.Fatalf("Error finding root directory: %v", err)
	}
	envPath := filepath.Join(rootDir, ".env")
	err = godotenv.Load(envPath)
	if err != nil {
		t.Fatalf("Error loading .env file from %s: %v", envPath, err)
	}

	// Get the API key.
	apiKey := os.Getenv("SPORTS_DATA_IO_KEY")
	if apiKey == "" {
		t.Fatal("SPORTS_DATA_IO_KEY is not set in .env")
	}

	// Create the SportsClient with the API key.
	sportsClient := SportsClient{APIKey: apiKey}

	// Define the season you want to query.
	season := "2025" // adjust as needed

	// Fetch all player season stats.
	stats, err := sportsClient.GetPlayerSeasonStats(season)
	if err != nil {
		t.Fatalf("Error fetching player season stats: %v", err)
	}

	// Filter the stats to find LeBron James.
	lebronStats := FindPlayerSeasonStats(stats, "LeBron", "James")
	if lebronStats == nil {
		t.Fatal("Could not find stats for LeBron James")
	}

	// Log the results for debugging.
	t.Logf("LeBron James Season Stats: %+v", lebronStats)
}
