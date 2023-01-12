package handlers

import (
	"DIA-NFT-Sales-Bot/config"
	"DIA-NFT-Sales-Bot/models"
	"database/sql"
	"fmt"
	"math/big"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

func AllSalesStopHandler(discordSession *discordgo.Session, interaction *discordgo.InteractionCreate) {
	var (
		message             string
		optionsMap          = ParseCommandOptions(interaction)
		subs                = models.Subscriptions{Command: "all_sales"}
		threshold, tExists  = optionsMap["threshold"]
		channel, chanExists = optionsMap["channel"]
		all, blockchain     = optionsMap["all"].BoolValue(), optionsMap["blockchain"].StringValue()
	)
	config.ActiveAllSalesMux.Lock()

	//Respond Channel is being Setup
	if all {
		message = "Deactivate AllSales Subscription for all Channels and Thresholds"

		go maps.Clear(config.ActiveAllSales)
		config.ActiveAllSalesKeys = []*big.Float{}
		go subs.UnsubscribeSalesUpdates()

	} else {
		if !tExists || !chanExists {
			message = "Invalid Threshold or Channel supplied"
		} else {
			channelID := channel.ChannelValue(discordSession).ID

			go unsubscribeAllSales(channelID, subs, threshold.StringValue(), blockchain)

			message = fmt.Sprintf("Deactivated AllSales Subscription for Threshold : %f ETH  on Channel: %s", threshold.FloatValue(), channel.Name)
		}
	}
	defer config.ActiveAllSalesMux.Unlock()

	err := discordSession.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})

	if err != nil {
		panic(err)
	}
}

func unsubscribeAllSales(channelID string, subs models.Subscriptions, threshold string, blockchain string) {
	subs.ChannelID, subs.Threshold, subs.Blockchain = sql.NullString{String: channelID, Valid: true}, threshold, blockchain

	go subs.UnsubscribeChannelSalesUpdates()
	thresholdbigInt := big.NewFloat(0)
	subscribedChannels := config.ActiveAllSales[thresholdbigInt][blockchain]
	for index, c := range subscribedChannels {
		if c == channelID {
			subscribedChannels = slices.Delete(subscribedChannels, index, index+1)
			config.ActiveAllSales[thresholdbigInt][blockchain] = subscribedChannels
			break
		}
	}
}
