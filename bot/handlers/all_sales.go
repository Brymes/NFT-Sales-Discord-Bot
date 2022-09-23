package handlers

import (
	"DIA-NFT-Sales-Bot/config"
	"DIA-NFT-Sales-Bot/models"
	"DIA-NFT-Sales-Bot/services"
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"sort"
)

func AllSalesHandler(discordSession *discordgo.Session, interaction *discordgo.InteractionCreate) {
	optionsMap := ParseCommandOptions(interaction)

	channel, threshold := optionsMap["channel"].ChannelValue(nil), optionsMap["threshold"].FloatValue()

	//Respond Channel is being Setup
	message := fmt.Sprintf("Setup Channel:%s \t to receive updates for Contract Address: %s", channel.Name, threshold)
	err := discordSession.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})

	if err != nil {
		panic(err)
	}

	//Check if WebSocket has started
	if !config.ActiveNftEventWS {
		services.StartEventWS()
	}

	//Add Details to Subscriptions in DB
	models.Subscriptions{
		Command:   "sales",
		ChannelID: sql.NullString{String: channel.ID, Valid: true},
		Threshold: threshold,
		Active:    true,
	}.SaveSubscription()

	//Add Details to AllSales[ChannelID]
	config.ActiveAllSalesMux.Lock()
	config.ActiveAllSales[threshold] = append(config.ActiveAllSales[threshold], channel.ID)
	config.ActiveAllSalesKeys = make([]float64, 0, len(config.ActiveAllSales[threshold]))

	for k := range config.ActiveAllSales {
		config.ActiveAllSalesKeys = append(config.ActiveAllSalesKeys, k)
		config.ActiveAllSalesReversed = append(config.ActiveAllSalesReversed, k)
	}

	sort.Float64s(config.ActiveAllSalesKeys)

	config.ActiveAllSalesMux.Unlock()

	//Follow Up has been Set up
	SendChannelSetupFollowUp("Channel setup complete & successful", discordSession, interaction)
}
