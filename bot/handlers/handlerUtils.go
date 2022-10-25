package handlers

import (
	"github.com/bwmarrin/discordgo"
)

var currencies = map[string]string{
	"ethereum": "ETH",
	"astar":    "SOL",
	"solana":   "ASTR",
}

func ParseCommandOptions(interaction *discordgo.InteractionCreate) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	// Access options in the order provided by the user.
	options := interaction.ApplicationCommandData().Options

	// Or convert the slice into a map
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	return optionMap
}

func SendChannelSetupFollowUp(message string, discordSession *discordgo.Session, interaction *discordgo.InteractionCreate) {
	_, err := discordSession.FollowupMessageCreate(interaction.Interaction, true, &discordgo.WebhookParams{Content: message})
	if err != nil {
		panic("Error Sending FollowUp Message " + err.Error())
	}
}

func SendComplexMessage(followUp discordgo.MessageSend, discordSession *discordgo.Session, interaction *discordgo.InteractionCreate) {
	_, err := discordSession.ChannelMessageSendComplex(interaction.ChannelID, &followUp)
	if err != nil {
		panic("Error Sending FollowUp Message " + err.Error())
	}
}
