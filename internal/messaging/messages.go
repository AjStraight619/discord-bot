package messaging

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func MessageAllMembers(dg *discordgo.Session, guildID, message string) {
	members, err := dg.GuildMembers(guildID, "", 1000)
	if err != nil {
		log.Printf("Error fetching members: %v", err)
		return
	}

	for _, member := range members {
		// Optionally skip bots
		if member.User.Bot {
			continue
		}
		SendDM(dg, member.User.ID, message)
	}
}

func SendDM(dg *discordgo.Session, userID, message string) {
	channel, err := dg.UserChannelCreate(userID)
	if err != nil {
		log.Printf("Error creating DM channel for user %s: %v", userID, err)
		return
	}

	_, err = dg.ChannelMessageSend(channel.ID, message)
	if err != nil {
		log.Printf("Error sending DM to user %s: %v", userID, err)
	} else {
		log.Printf("DM sent to user %s", userID)
	}
}

func ListGuildMembersByName(dg *discordgo.Session, guildName string) {
	var targetGuild *discordgo.Guild
	// Loop over the cached guilds to find the one matching the name.
	for _, guild := range dg.State.Guilds {
		if guild.Name == guildName {
			targetGuild = guild
			break
		}
	}
	if targetGuild == nil {
		log.Printf("Guild '%s' not found in session state.", guildName)
		return
	}

	log.Printf("Listing members of guild '%s' (ID: %s):", targetGuild.Name, targetGuild.ID)
	// Loop over the cached members of the guild.
	for _, member := range targetGuild.Members {
		log.Printf("Member: %s (ID: %s)", member.User.Username, member.User.ID)
	}
}
