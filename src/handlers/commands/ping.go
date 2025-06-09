package commands

import "github.com/bwmarrin/discordgo"

func PingPong(s *discordgo.Session, m *discordgo.InteractionCreate) {
	s.InteractionRespond(m.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong!",
		},
	})
}
