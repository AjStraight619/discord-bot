package members

import "github.com/bwmarrin/discordgo"

func FetchAllGuildMembers(dg *discordgo.Session, guildID string) ([]*discordgo.Member, error) {
	var allMembers []*discordgo.Member
	lastMemberID := ""

	for {
		// Fetch up to 1000 members at a time.
		members, err := dg.GuildMembers(guildID, lastMemberID, 1000)
		if err != nil {
			return allMembers, err
		}

		// If no members were returned, we have fetched all.
		if len(members) == 0 {
			break
		}

		allMembers = append(allMembers, members...)
		lastMemberID = members[len(members)-1].User.ID

		if len(members) < 1000 {
			break
		}
	}
	return allMembers, nil
}
