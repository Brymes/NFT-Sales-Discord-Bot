package utils

import (
	"DIA-NFT-Sales-Bot/config"
	log "DIA-NFT-Sales-Bot/debug"
	"fmt"

	"runtime/debug"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/olekukonko/tablewriter"
)

func sendPanicMessage(panicTable [][]string, discordSession *discordgo.Session) {
	str := &strings.Builder{}
	table := tablewriter.NewWriter(str)
	table.SetHeader([]string{"Command", "Desc"})
	table.AppendBulk(panicTable)
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetHeaderLine(true)
	table.Render()
	msg := "```" + str.String() + "```"

	if len(msg) > 1950 {
		msg = msg[:1950]
	}

	_, err := discordSession.ChannelMessageSend(config.PanicChannelID, msg)
	if err != nil {
		log.Log.Println(err)
		log.Log.Println("Error Sending Panic message")
	}

}

func HandlePanic(discordSession *discordgo.Session, customMessage string) {
	var msg = [][]string{
		{"message", customMessage},
	}
	//Notify Admin of any uncaught errors
	if err := recover(); err != nil {
		stack := string(debug.Stack())

		log.Log.Println(err)
		log.Log.Println(customMessage)
		log.Log.Println(stack)

		msg = append(msg, []string{"Error", fmt.Sprintf("%v", err)})
		msg = append(msg, []string{"Call Stack", stack})

		go sendPanicMessage(msg, discordSession)
	}
}
