package main

import (
	"DIA-NFT-Sales-Bot/bot"
	"DIA-NFT-Sales-Bot/config"
	"DIA-NFT-Sales-Bot/models"
	"DIA-NFT-Sales-Bot/services"
	"DIA-NFT-Sales-Bot/utils"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	config.InitPanicChannel()
	config.InitDb()
	models.InitMigrations()
	startWS := models.LoadCurrentSubscriptions()
	bot.InitBot()
	if startWS {
		services.StartEventWS()
	}
}

func main() {

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session .
	defer config.DiscordBot.Close()
	defer bot.DeRegisterCommands(config.DiscordBot)
	defer config.ShutDownWS()

	defer utils.HandlePanic(config.DiscordBot, "")
}
