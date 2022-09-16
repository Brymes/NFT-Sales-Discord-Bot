package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
)

func InitBot(Token string) *discordgo.Session {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return nil
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

	return dg

}
