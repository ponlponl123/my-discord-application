package utils

import (
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func SetDefaultActivity(s *discordgo.Session) {
	total_guilds := len(s.State.Guilds)
	if total_guilds <= 1 {
		log.Println("There is not have much guilds to show")
		return
	}
	err := s.UpdateStatusComplex(discordgo.UpdateStatusData{
		Status: string(discordgo.StatusOnline),
		Activities: []*discordgo.Activity{
			{
				Name: strconv.Itoa(total_guilds) + " servers",
				Type: discordgo.ActivityTypeWatching,
			},
		},
	})
	if err != nil {
		log.Println("Failed to set default activity")
		log.Println("Error:", err)
	}
}
