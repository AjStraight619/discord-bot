package messaging

import (
	"fmt"
	"log"
	"strings"

	"github.com/AjStraight619/discord-bot/internal/apiclients"
	"github.com/AjStraight619/discord-bot/internal/members"
	"github.com/AjStraight619/discord-bot/internal/models"
	sportsutils "github.com/AjStraight619/discord-bot/internal/sports_utils"
	"github.com/AjStraight619/discord-bot/internal/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
)

// CronMessages holds the data for scheduling messages.
type CronMessage struct {
	Session  *discordgo.Session
	duration int
	cron     *cron.Cron
	cronMsg  string
	schedule string
}

// InitCron initializes a new CronMessages instance.
func InitCron(dg *discordgo.Session, dur int) *CronMessage {
	return &CronMessage{
		Session:  dg,
		duration: dur,
		cron:     cron.New(),
		cronMsg:  "",
		schedule: "",
	}
}

// StartCydCron loads team stats, formats the message, and schedules a DM cron job.
func (cm *CronMessage) StartCydCron() {
	// Load team data.

	// TODO: Move LoadTeams to sports utils package. Bot should not know about loading teams.
	teams := sportsutils.LoadTeams()
	if teams == nil {
		log.Printf("Failed to load teams")
		return
	}

	lakers, err := teams.FindTeam("Lakers")
	if err != nil {
		log.Printf("Failed to find team: %v", err)
		return
	}

	teamStats, err := apiclients.GetTeamStatistics(lakers.ID, "2024", "REG")
	if err != nil {
		log.Printf("Failed to get team stats: %v", err)
		return
	}

	log.Printf("Team Stats: %+v", teamStats)
	log.Printf("Team: %s", teamStats.Name)
	log.Printf("Market: %s", teamStats.Market)
	log.Printf("Season: %d %s", teamStats.Season.Year, teamStats.Season.Type)
	log.Printf("Games Played: %d", teamStats.OwnRecord.Total.GamesPlayed)

	// Look for LeBron James in the players list and format his stats.
	for _, player := range teamStats.Players {
		if strings.EqualFold("LeBron James", player.FullName) {
			log.Printf("Found player: %s", player.FullName)
			// Optionally log some averages:
			log.Printf("Minutes: %.2f | Points: %.2f", player.Averages.Minutes, player.Averages.Points)
			cm.cronMsg = FormatPlayerStatsMessage(player)
			break
		}
	}

	// Find the guild by name.
	guild := utils.FindGuildByName(cm.Session, "King's Landing")
	if guild == nil {
		log.Printf("Error finding guild...")
		return
	}

	// Fetch all guild members.
	allMembers, err := members.FetchAllGuildMembers(cm.Session, guild.ID)
	if err != nil {
		log.Printf("Error fetching guild members: %v", err)
		return
	}

	// Find the target member with username "cydstynine".
	var targetUserID string
	for _, member := range allMembers {
		if strings.EqualFold(member.User.Username, "cydstynine") {
			targetUserID = member.User.ID
			break
		}
	}
	if targetUserID == "" {
		log.Printf("Member with username 'cydstynine' not found")
		return
	}

	schedule := "@daily"

	// Add the cron job to send a DM with the formatted message.
	_, err = cm.cron.AddFunc(schedule, func() {
		// Use your SendDM function (from your bot package or wherever it is defined).
		SendDM(cm.Session, targetUserID, cm.cronMsg)
	})
	if err != nil {
		log.Printf("Error scheduling cron job: %v", err)
		return
	}

	// Start the cron scheduler.
	cm.cron.Start()
	log.Printf("Cron job started for user %s with schedule: %s", targetUserID, schedule)
}

// FormatPlayerStatsMessage formats the player's averages into a message string.
func FormatPlayerStatsMessage(player models.Player) string {
	return fmt.Sprintf(
		"\n\nHere are your average stats for %s:\n"+
			"Minutes: %.2f\n"+
			"Points: %.2f\n"+
			"Offensive Rebounds: %.2f\n"+
			"Defensive Rebounds: %.2f\n"+
			"Rebounds: %.2f\n"+
			"Assists: %.2f\n"+
			"Steals: %.2f\n"+
			"Blocks: %.2f\n"+
			"Turnovers: %.2f\n"+
			"Personal Fouls: %.2f\n",
		player.FullName,
		player.Averages.Minutes,
		player.Averages.Points,
		player.Averages.OffRebounds,
		player.Averages.DefRebounds,
		player.Averages.Rebounds,
		player.Averages.Assists,
		player.Averages.Steals,
		player.Averages.Blocks,
		player.Averages.Turnovers,
		player.Averages.PersonalFouls,
	)
}
