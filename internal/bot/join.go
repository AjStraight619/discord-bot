package bot

import (
	"github.com/bwmarrin/discordgo"
)

type JoinCommand struct{}

func (jc JoinCommand) Execute(b *BotController, msg *discordgo.MessageCreate, options []string) {
	b.joinUserChannel(msg.GuildID, msg.Author.ID, true, true)
}

func (jc JoinCommand) Help() string {
	return "!join - Join the voice channel of the user who sent the command."
}
