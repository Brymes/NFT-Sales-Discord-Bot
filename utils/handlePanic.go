package utils

import (
	"DIA-NFT-Sales-Bot/config"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/olekukonko/tablewriter"
	"log"
	"runtime/debug"
	"strings"
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
		log.Println(err)
		log.Println("Error Sending Panic message")
	}

}

func HandlePanic(discordSession *discordgo.Session, customMessage string) {
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

		go sendPanicMessage(msg, discordSession)
	}
}
