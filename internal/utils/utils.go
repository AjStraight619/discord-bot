package utils

import "github.com/bwmarrin/discordgo"

func FindGuildByName(dg *discordgo.Session, guildName string) *discordgo.Guild {
	for _, guild := range dg.State.Guilds {
		if guild.Name == guildName {
			return guild
		}
	}
	return nil
}
