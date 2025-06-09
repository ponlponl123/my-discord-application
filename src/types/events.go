package types

import (
	"my-discord-bot/src/handlers/commands"
	"my-discord-bot/src/handlers/events"

	"github.com/bwmarrin/discordgo"
)

// Commands stores all application commands
var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "ping",
		Description: "Ping Pong!",
	},
	{
		Name:        "referal",
		Description: "Grant you a specific permission by referal code.",
	},
}

// RegisteredCommands stores all registered application commands
var RegisteredCommands = make([]*discordgo.ApplicationCommand, 0)

// CommandHandlers maps command names to their handler functions
var CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"ping":    commands.PingPong,
	"referal": commands.Referal,
}

// ModalHandlers maps modal IDs to their handler functions
var ModalHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"referal": commands.ReferalModal,
}

// UniversalEvents is a slice of functions to register event handlers
var UniversalEvents = []func(s *discordgo.Session){
	func(s *discordgo.Session) {
		s.AddHandler(events.MemberJoin)
	},
	func(s *discordgo.Session) {
		s.AddHandler(events.MemberLeave)
	},
	func(s *discordgo.Session) {
		s.AddHandler(events.MessageCreate)
	},
	func(s *discordgo.Session) {
		s.AddHandler(events.MessageDelete)
	},
	func(s *discordgo.Session) {
		s.AddHandler(events.MessageUpdate)
	},
	func(s *discordgo.Session) {
		s.AddHandler(events.MessageReactionAdd)
	},
}
