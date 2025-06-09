package handlers

import (
	"log"
	"my-discord-bot/src/types"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// RegisterCommands sets up all bot command handlers
func RegisterCommands(discord *discordgo.Session) {
	// Add message handler for all incoming messages
	discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := types.CommandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionModalSubmit:
			if h, ok := types.ModalHandlers[strings.Split(i.ModalSubmitData().CustomID, "_")[0]]; ok {
				h(s, i)
			}
		}
	})

	if len(types.Commands) == 0 {
		log.Println("There are no commands to register.")
		return
	}
	if len(discord.State.Guilds) == 0 {
		log.Println("There is no guilds to register commands for.")
		return
	}
	for _, g := range discord.State.Guilds {
		for _, v := range types.Commands {
			cmd, err := discord.ApplicationCommandCreate(discord.State.User.ID, g.ID, v)
			if err != nil {
				log.Panicf("Cannot create '%v' command: %v\n", v.Name, err)
			}
			log.Printf("Registered %v command for %v\n", v.Name, g.ID)
			types.RegisteredCommands = append(types.RegisteredCommands, cmd)
		}
	}
}

func CleanUpCommands(discord *discordgo.Session) {
	if len(types.RegisteredCommands) == 0 {
		log.Println("There are no registered commands to clean up.")
		return
	}
	for _, v := range types.RegisteredCommands {
		err := discord.ApplicationCommandDelete(discord.State.User.ID, v.GuildID, v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v\n", v.Name, err)
		}
		log.Printf("Deleted command %v for %v\n", v.Name, v.GuildID)
	}
}
