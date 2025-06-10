package types

import "github.com/bwmarrin/discordgo"

// Commands stores all application commands
var Commands = []*discordgo.ApplicationCommand{}

// RegisteredCommands stores all registered application commands
var RegisteredCommands = make([]*discordgo.ApplicationCommand, 0)

// CommandHandlers maps command names to their handler functions
var CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
