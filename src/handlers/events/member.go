package events

import (
	"my-discord-bot/src/handlers/extra"

	"github.com/bwmarrin/discordgo"
)

func MemberJoin(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	extra.Referal_MemberAddedEvent(s, m)
}

func MemberLeave(s *discordgo.Session, m *discordgo.GuildMemberRemove) {

}
