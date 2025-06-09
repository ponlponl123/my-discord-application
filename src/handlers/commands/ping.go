package commands

import "github.com/bwmarrin/discordgo"

func PingPong(s *discordgo.Session, m *discordgo.InteractionCreate) {
	s.ChannelMessageSend(m.ChannelID, "Pong!")
}
