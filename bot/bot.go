package bot

import (
	"DIA-NFT-Sales-Bot/config"
	log "DIA-NFT-Sales-Bot/debug"
	"os"

	"github.com/bwmarrin/discordgo"
)

func InitBot() {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	token = "MTAxOTI3MDc1ODE4ODQ1ODA4NQ.GrmnqE._4fMGh82Dux7vpoHT5A7VxwbbMZmcpJrciFvgA"
	if token == "" {
		log.Log.Fatalln("DISCORD_BOT_TOKEN environment variable not set")
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Log.Fatalln("error creating Discord session,", err)
	}

	//Ensure Messages only come in from Guilds
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	RegisterHandlers(dg)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Log.Println("error opening connection,", err)
		log.Log.Fatal("Error Initializing bot")
	}

	RegisterCommands(dg)

	config.DiscordBot = dg

}
