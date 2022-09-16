package bot

import (
	"DIA-NFT-Sales-Bot/bot/handlers"
	"github.com/bwmarrin/discordgo"
)

var (
	SlashCommands = map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
		"help": handlers.HelpHandler,
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

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	switch m.Content {
	case "ping":
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	case "pong":
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	default:
		handlers.SendHelpText(s, m)
	}
}

func slashCommandHandler(discordSession *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if handler, ok := SlashCommands[interaction.ApplicationCommandData().Name]; ok {
		handler(discordSession, interaction)
	} else {
		SlashCommands["help"](discordSession, interaction)
	}
}
