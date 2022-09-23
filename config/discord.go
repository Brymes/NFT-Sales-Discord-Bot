package config

import (
	"DIA-NFT-Sales-Bot/bot"
	"log"
	"os"
)

func InitializeBotSession() {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatalln("Bot Token environment variable not set")
	}

	DiscordBot = bot.InitBot(token)
	return
}
