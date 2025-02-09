package models

// SRTeam represents the team data returned from the API.
type SRTeam struct {
	ID        string        `json:"id"`
	Name      string        `json:"name"`
	Market    string        `json:"market"`
	SRID      string        `json:"sr_id"`
	Reference string        `json:"reference"`
	Season    Season        `json:"season"`
	OwnRecord TeamOwnRecord `json:"own_record"`
	Players   []Player      `json:"players"`
}

// Season holds the season information for a team.
type Season struct {
	ID   string `json:"id"`
	Year int    `json:"year"`
	Type string `json:"type"`
}

// TeamOwnRecord contains both total and average stats for a team.
type TeamOwnRecord struct {
	Total   TeamTotalStats   `json:"total"`
	Average TeamAverageStats `json:"average"` // Optional: add if you want averages too.
}

// TeamTotalStats contains some of the total statistics for a team.
// For now, we only grab games_played, but you can add more fields as required.
type TeamTotalStats struct {
	GamesPlayed int `json:"games_played"`
	// Add additional fields if needed.
}

// TeamAverageStats contains average statistics for a team.
// You can fill this out with more fields if needed.
type TeamAverageStats struct {
	// Example:
	// Points           float64 `json:"points"`
	// Rebounds         float64 `json:"rebounds"`
}

// Player represents a player returned from the API.
type Player struct {
	ID              string         `json:"id"`
	FullName        string         `json:"full_name"`
	FirstName       string         `json:"first_name"`
	LastName        string         `json:"last_name"`
	Position        string         `json:"position"`
	PrimaryPosition string         `json:"primary_position"`
	JerseyNumber    string         `json:"jersey_number"`
	SRID            string         `json:"sr_id"`
	Reference       string         `json:"reference"`
	Averages        PlayerAverages `json:"average"`
	// You can also add a field for totals if needed.
}

// PlayerAverages holds all the averaged statistics for a player.
type PlayerAverages struct {
	Minutes            float64 `json:"minutes"`
	Points             float64 `json:"points"`
	OffRebounds        float64 `json:"off_rebounds"`
	DefRebounds        float64 `json:"def_rebounds"`
	Rebounds           float64 `json:"rebounds"`
	Assists            float64 `json:"assists"`
	Steals             float64 `json:"steals"`
	Blocks             float64 `json:"blocks"`
	Turnovers          float64 `json:"turnovers"`
	PersonalFouls      float64 `json:"personal_fouls"`
	FlagrantFouls      float64 `json:"flagrant_fouls"`
	BlockedAtt         float64 `json:"blocked_att"`
	FieldGoalsMade     float64 `json:"field_goals_made"`
	FieldGoalsAtt      float64 `json:"field_goals_att"`
	ThreePointsMade    float64 `json:"three_points_made"`
	ThreePointsAtt     float64 `json:"three_points_att"`
	FreeThrowsMade     float64 `json:"free_throws_made"`
	FreeThrowsAtt      float64 `json:"free_throws_att"`
	TwoPointsMade      float64 `json:"two_points_made"`
	TwoPointsAtt       float64 `json:"two_points_att"`
	Efficiency         float64 `json:"efficiency"`
	TrueShootingAtt    float64 `json:"true_shooting_att"`
	PointsOffTurnovers float64 `json:"points_off_turnovers"`
	PointsInPaintMade  float64 `json:"points_in_paint_made"`
	PointsInPaintAtt   float64 `json:"points_in_paint_att"`
	PointsInPaint      float64 `json:"points_in_paint"`
	FoulsDrawn         float64 `json:"fouls_drawn"`
	OffensiveFouls     float64 `json:"offensive_fouls"`
	FastBreakPts       float64 `json:"fast_break_pts"`
	FastBreakAtt       float64 `json:"fast_break_att"`
	FastBreakMade      float64 `json:"fast_break_made"`
	SecondChancePts    float64 `json:"second_chance_pts"`
	SecondChanceAtt    float64 `json:"second_chance_att"`
	SecondChanceMade   float64 `json:"second_chance_made"`
}
