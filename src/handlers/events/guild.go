package events

import (
	"my-discord-bot/src/utils"

	"github.com/bwmarrin/discordgo"
)

func GuildCreate(s *discordgo.Session, m *discordgo.GuildCreate) {
	utils.SetDefaultActivity(s)
}

func GuildDelete(s *discordgo.Session, m *discordgo.GuildDelete) {
	utils.SetDefaultActivity(s)
}
