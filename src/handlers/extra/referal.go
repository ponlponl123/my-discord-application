package extra

import (
	"github.com/bwmarrin/discordgo"
)

func Referal_MemberAddedEvent(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	_, err := s.Guild(m.GuildID)
	if err != nil {
		return
	}
}
