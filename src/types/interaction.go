package types

import (
	"github.com/bwmarrin/discordgo"
)

// ModalHandlers maps modal IDs to their handler functions
var ModalHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
