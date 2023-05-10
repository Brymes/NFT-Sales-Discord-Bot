package bot

import (
	"DIA-NFT-Sales-Bot/config"
	log "DIA-NFT-Sales-Bot/debug"
	"os"

	"github.com/bwmarrin/discordgo"
)

func InitBot() {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatalln("DISCORD_BOT_TOKEN environment variable not set")
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalln("error creating Discord session,", err)
	}

	//Ensure Messages only come in from Guilds
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	RegisterHandlers(dg)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Println("error opening connection,", err)
		log.Fatal("Error Initializing bot")
	}

	RegisterCommands(dg)

	config.DiscordBot = dg

}
