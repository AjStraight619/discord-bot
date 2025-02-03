package bot

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (b *BotController) JoinChannel(guildID string, userID string) *discordgo.VoiceConnection {
	guild, err := b.Session.Guild(guildID)

	if err != nil {
		log.Println("Error fetching guild:", err)
		return nil
	}

	var voiceChannelID string
	for _, vs := range guild.VoiceStates {
		if vs.UserID == userID {
			voiceChannelID = vs.ChannelID
			break
		}
	}

	if voiceChannelID == "" {
		b.Session.ChannelMessageSend(guildID, "⚠ You must be in a voice channel for me to join!")
		return nil
	}

	vc, err := b.Session.ChannelVoiceJoin(guildID, voiceChannelID, false, true)
	if err != nil {
		log.Println("❌ Error joining voice channel:", err)
		return nil
	}

	b.isBotInChannel = true
	return vc
}

func (b *BotController) JoinChannelByName(guildID, channelName string, deafened bool) *discordgo.VoiceConnection {
	// ✅ Get all channels in the server
	channels, err := b.Session.GuildChannels(guildID)
	if err != nil {
		log.Println("❌ Error fetching channels:", err)
		return nil
	}

	var voiceChannelID string

	// ✅ Search for a voice channel that matches the name
	for _, channel := range channels {
		if strings.EqualFold(channel.Name, channelName) && channel.Type == discordgo.ChannelTypeGuildVoice {
			voiceChannelID = channel.ID
			break
		}
	}

	// ✅ If no matching channel is found, send an error message
	if voiceChannelID == "" {
		b.Session.ChannelMessageSend(guildID, fmt.Sprintf("⚠ No voice channel named **%s** found!", channelName))
		return nil
	}

	// ✅ Join the detected voice channel
	vc, err := b.Session.ChannelVoiceJoin(guildID, voiceChannelID, false, deafened)
	if err != nil {
		log.Println("❌ Error joining voice channel:", err)
		return nil
	}

	log.Println("✅ Joined voice channel:", channelName)
	return vc

}
