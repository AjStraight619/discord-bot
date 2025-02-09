package bot

import "log"

func (b *BotController) ListGuildChannels() {

	log.Printf("Length of guilds: %d", len(b.Session.State.Guilds))

	for _, guild := range b.Session.State.Guilds {
		channels, err := b.Session.GuildChannels(guild.ID)
		if err != nil {
			log.Printf("Error retrieving channels for guild %s: %v", guild.ID, err)
			continue
		}
		log.Printf("Channels in %s:", guild.Name)
		for _, channel := range channels {
			log.Printf(" - %s (ID: %s) [Type: %d]", channel.Name, channel.ID, channel.Type)
		}
	}
}
