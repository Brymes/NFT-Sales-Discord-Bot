package handlers

import (
	"DIA-NFT-Sales-Bot/config"
	"DIA-NFT-Sales-Bot/models"
	"DIA-NFT-Sales-Bot/services"
	"DIA-NFT-Sales-Bot/utils"
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"sort"
	"strings"
)

func AllSalesHandler(discordSession *discordgo.Session, interaction *discordgo.InteractionCreate) {
	optionsMap := ParseCommandOptions(interaction)

	channel, threshold, blockchain := optionsMap["channel"].ChannelValue(discordSession), optionsMap["threshold"].FloatValue(), optionsMap["blockchain"].StringValue()

	//Respond Channel is being Setup
	message := fmt.Sprintf("Setup Channel : %s  to receive updates for Sales above price threshold : %f %s", channel.Name, threshold, currencies[strings.ToLower(blockchain)])
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
		Command:    "all_sales",
		Blockchain: blockchain,
		ChannelID:  sql.NullString{String: channel.ID, Valid: true},
		Threshold:  threshold,
		Active:     true,
	}.SaveSubscription()

	addDetailsToMap(threshold, channel.ID, blockchain)

	//Follow Up has been Set up
	SendChannelSetupFollowUp("Channel setup complete & successful", discordSession, interaction)
}

func addDetailsToMap(threshold float64, channelID string, blockchain string) {

	//Add Details to AllSales[ChannelID]
	config.ActiveAllSalesMux.Lock()
	data := config.ActiveAllSales[threshold]
	if len(data) == 0 {
		data = map[string][]string{}
	}
	subscribedChannels := data[blockchain]
	subscribedChannels = append(subscribedChannels, channelID)

	data[blockchain] = utils.RemoveArrayDuplicates(subscribedChannels)
	config.ActiveAllSales[threshold] = data

	config.ActiveAllSalesKeys = make([]float64, 0, len(subscribedChannels))

	for k := range config.ActiveAllSales {
		config.ActiveAllSalesKeys = append(config.ActiveAllSalesKeys, k)
	}

	sort.Float64s(config.ActiveAllSalesKeys)

	defer config.ActiveAllSalesMux.Unlock()
}
