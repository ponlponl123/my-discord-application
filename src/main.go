package main

import (
	"log"
	"my-discord-bot/src/handlers"
	"my-discord-bot/src/types"
	"my-discord-bot/src/utils"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

func main() {
	log.Println("Starting discord bot...")

	token := utils.GetEnv("DISCORD_TOKEN", "")
	if token == "" {
		log.Println("No token provided, exiting...")
		return
	}
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Println("Error creating Discord session")
		log.Fatalln(err)
	}
	log.Println("Discord session created")
	discord.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up!, logged in as ", r.User.Username, "#", r.User.Discriminator, "(", r.User.ID, ")")
		utils.SetDefaultActivity(discord)
		utils.DB, err = utils.ConnectDB()
		if err != nil {
			log.Println("Error connecting to database")
			log.Fatalln(err)
		}
	})

	// Open connection to Discord
	err = discord.Open()
	if err != nil {
		log.Println("Open Discord connection failed")
		log.Fatalln(err)
	}
	defer discord.Close()

	// Register all commandsm events
	handlers.RegisterCommands(discord)
	for _, v := range types.UniversalEvents {
		v(discord)
	}

	stchan := make(chan os.Signal, 1)
	signal.Notify(stchan, syscall.SIGTERM, os.Interrupt, syscall.SIGSEGV)
end:
	for {
		select {
		case <-stchan:
			log.Println("Interrupt signal received, shutting down...")
			log.Println("Cleaning up registered commands before exit...")
			handlers.CleanUpCommands(discord)
			log.Println("Closing database connection...")
			utils.DB.Close()
			log.Println("Exiting...")
			break end
		default:
		}
		time.Sleep(time.Second)
	}
}
