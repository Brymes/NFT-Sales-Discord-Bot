package config

import (
	"DIA-NFT-Sales-Bot/bot"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
)

func InitializeBotSession() *discordgo.Session {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatalln("Bot Token environment variable not set")
	}

	session := bot.InitBot(token)

	return session
}
