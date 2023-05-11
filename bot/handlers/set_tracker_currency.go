package handlers

import (
	"DIA-NFT-Sales-Bot/config"
	"DIA-NFT-Sales-Bot/models"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func SetCurrencyHandler(discordSession *discordgo.Session, interaction *discordgo.InteractionCreate) {
	optionsMap := ParseCommandOptions(interaction)

	currency := optionsMap["currency"].StringValue()

	//Respond Channel is being Setup
	message := fmt.Sprintf("Set Floor price tracker to Currency: \"%s\"", currency)
	err := discordSession.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})

	if err != nil {
		panic(err)
	}

	config.TrackerCurrency = currency
	cm := models.ConfigModel{TrackerCurrency: currency}
	cm.SaveConfig()
}
