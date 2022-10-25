package handlers

import (
	"DIA-NFT-Sales-Bot/config"
	"DIA-NFT-Sales-Bot/models"
	"DIA-NFT-Sales-Bot/services"
	"DIA-NFT-Sales-Bot/utils"
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func SalesHandler(discordSession *discordgo.Session, interaction *discordgo.InteractionCreate) {
	optionsMap := ParseCommandOptions(interaction)

	channel, address, blockchain := optionsMap["channel"].ChannelValue(discordSession), optionsMap["contract_address"].StringValue(), optionsMap["blockchain"].StringValue()

	//Respond Channel is being Setup
	message := fmt.Sprintf("Setup Channel:%s \t to receive sales updates for Contract Address: %s", channel.Name, address)
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
		Command:    "sales",
		Blockchain: blockchain,
		ChannelID:  sql.NullString{String: channel.ID, Valid: true},
		Address:    sql.NullString{String: address, Valid: true},
		Active:     true,
	}.SaveSubscription()

	//Add Details to AllSales[ChannelID]
	config.ActiveSalesMux.Lock()

	subscribedChannels := config.ActiveSales[address][blockchain]
	subscribedChannels = append(subscribedChannels, channel.ID)
	config.ActiveSales[address][blockchain] = utils.RemoveArrayDuplicates(subscribedChannels)

	defer config.ActiveSalesMux.Unlock()

	//Follow Up has been Set up
	SendChannelSetupFollowUp("Channel setup complete & successful", discordSession, interaction)
}
