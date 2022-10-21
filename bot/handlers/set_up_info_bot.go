package handlers

import (
	"DIA-NFT-Sales-Bot/config"
	"DIA-NFT-Sales-Bot/models"
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

func SetUpInfoBotHandler(discordSession *discordgo.Session, interaction *discordgo.InteractionCreate) {
	optionsMap := ParseCommandOptions(interaction)

	channel, address, blockchain := optionsMap["channel"].ChannelValue(discordSession), optionsMap["contract_address"].StringValue(), optionsMap["blockchain"].StringValue()

	//Respond Channel is being Setup
	message := fmt.Sprintf("Enable Channel:%s \t to accept commands volume, floor_price, last_trades. for Contract Address: %s", channel.Name, address)
	err := discordSession.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})

	if err != nil {
		panic(err)
	}

	address = strings.ToUpper(address)
	//Add Details to Subscriptions in DB
	models.Subscriptions{
		Command:    "sales",
		Blockchain: blockchain,
		ChannelID:  sql.NullString{String: channel.ID, Valid: true},
		Address:    sql.NullString{String: address, Valid: true},
		Active:     true,
	}.SaveSubscription()

	//Add Details to AllSales[ChannelID]
	config.ActiveSalesInfoMux.Lock()

	_, found := config.ActiveSalesInfoBot[channel.ID]
	if found {
		SendChannelSetupFollowUp("Channel setup failed. Cannot set up info bot on same channel with multiple addresses", discordSession, interaction)
		return
	} else {
		config.ActiveSalesInfoBot[channel.ID] = map[string]string{
			"address": address, "blockchain": blockchain,
		}
	}

	defer config.ActiveSalesInfoMux.Unlock()

	//Follow Up has been Set up
	SendChannelSetupFollowUp("Channel setup for commands volume, floor_price, last_trades complete & successful", discordSession, interaction)
}
