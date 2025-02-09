package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type LeaveCommand struct{}

func (lc LeaveCommand) Execute(b *BotController, msg *discordgo.MessageCreate, options []string) {
	b.LeaveVoiceChannel()
}

func (lc LeaveCommand) Help() string {
	return "!leave - leave voice channel"
}

func (b *BotController) LeaveVoiceChannel() {
	if b.VoiceConn != nil {
		log.Println("👋 Leaving voice channel...")
		if b.VoiceConn.OpusRecv != nil {
			close(b.VoiceConn.OpusRecv)
		}
		b.VoiceConn.Disconnect()
		b.VoiceConn = nil
		b.isBotInChannel = false
		b.Session.ChannelMessageSend(b.VoiceTextChannelID, "✅ Left the voice channel due to inactivity.")
	}
}
