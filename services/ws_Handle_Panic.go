package services

import (
	"DIA-NFT-Sales-Bot/config"
	"DIA-NFT-Sales-Bot/utils"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"runtime/debug"
	"time"
)

func WSHandlePanic(discordSession *discordgo.Session, customMessage string, logger *log.Logger) {
	defer utils.HandlePanic(discordSession, customMessage)
	var msg = [][]string{
		{"message", customMessage},
	}
	//Notify Admin of any uncaught errors
	if err := recover(); err != nil {
		stack := string(debug.Stack())

		log.Println(err)
		log.Println(customMessage)
		log.Println(stack)

		msg = append(msg, []string{"Error", fmt.Sprintf("%v", err)})
		msg = append(msg, []string{"Call Stack", stack})

		log.Println("restarting Connection")
		config.ShutDownWS()
		time.Sleep(5 * time.Second)
		go ConnectToService(logger)
		config.ActiveNftEventWS = true
		logger.Println("Restarted WebSocket")
	}
}
