package types

import (
	"my-discord-bot/src/handlers/commands"

	"github.com/bwmarrin/discordgo"
)

// Commands stores all application commands
var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "ping",
		Description: "Ping Pong!",
	},
}

// RegisteredCommands stores all registered application commands
var RegisteredCommands = make([]*discordgo.ApplicationCommand, 0)

// CommandHandlers maps command names to their handler functions
var CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"ping": commands.PingPong,
}
