package members

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func ListGuildMembersByName(dg *discordgo.Session, guildName string) []*discordgo.Member {
	var targetGuild *discordgo.Guild
	var allMembers []*discordgo.Member
	// Loop over the cached guilds to find the one matching the name.
	for _, guild := range dg.State.Guilds {
		if guild.Name == guildName {
			targetGuild = guild
			break
		}
	}
	if targetGuild == nil {
		log.Printf("Guild '%s' not found in session state.", guildName)
		return nil
	}

	log.Printf("Listing members of guild '%s' (ID: %s):", targetGuild.Name, targetGuild.ID)
	// Loop over the cached members of the guild.
	for _, member := range targetGuild.Members {
		log.Printf("Member: %s (ID: %s)", member.User.Username, member.User.ID)
		allMembers = append(allMembers, member)
	}

	return allMembers
}
