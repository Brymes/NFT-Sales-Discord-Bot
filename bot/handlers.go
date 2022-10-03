package bot

import (
	"DIA-NFT-Sales-Bot/bot/handlers"
	"DIA-NFT-Sales-Bot/utils"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

var (
	SlashCommands = map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
		"help":           handlers.HelpHandler,
		"subscriptions":  handlers.SubscriptionsHandler,
		"sales":          handlers.SalesHandler,
		"sales_stop":     handlers.SalesStopHandler,
		"floor":          handlers.FloorHandler,
		"all_sales":      handlers.AllSalesHandler,
		"all_sales_stop": handlers.AllSalesStopHandler,
		"stop_all":       handlers.StopAllHandler,
	}
)

func RegisterHandlers(discordSession *discordgo.Session) {
	// Register the slash commands Handler func as a callback for MessageCreate events.
	discordSession.AddHandler(slashCommandHandler)

	// Register the message Handler func as a callback for MessageCreate events.
	discordSession.AddHandler(messageHandler)
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	defer utils.HandlePanic(s, fmt.Sprintf("Error Handling message %v", m.Content))

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	switch m.Content {
	case "ping":
		_, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
		if err != nil {
			panic(err)
		}
	case "pong":
		_, err := s.ChannelMessageSend(m.ChannelID, "Ping!")
		if err != nil {
			panic(err)
		}
	default:
		handlers.SendHelpText(s, m)
	}
}

func slashCommandHandler(discordSession *discordgo.Session, interaction *discordgo.InteractionCreate) {
	defer utils.HandlePanic(discordSession, fmt.Sprintf("Error Handling command %s", interaction.ApplicationCommandData().Name))

	if handler, ok := SlashCommands[interaction.ApplicationCommandData().Name]; ok {
		handler(discordSession, interaction)
	} else {
		SlashCommands["help"](discordSession, interaction)
	}
}
