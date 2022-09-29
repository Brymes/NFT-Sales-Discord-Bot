package handlers

import (
	"DIA-NFT-Sales-Bot/config"
	"DIA-NFT-Sales-Bot/models"
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

func AllSalesStopHandler(discordSession *discordgo.Session, interaction *discordgo.InteractionCreate) {
	var (
		message                 string
		optionsMap              = ParseCommandOptions(interaction)
		subs                    = models.Subscriptions{Command: "all_sales"}
		threshold, channel, all = optionsMap["threshold"].FloatValue(), optionsMap["channel"].ChannelValue(nil), optionsMap["all"].BoolValue()
	)
	config.ActiveAllSalesMux.Lock()

	//Respond Channel is being Setup
	if all {
		message = "Deactivate AllSales Subscription for all Channels and Thresholds"

		go maps.Clear(config.ActiveAllSales)
		go subs.UnsubscribeSalesUpdates()

	} else {
		subs.ChannelID, subs.Threshold = sql.NullString{String: channel.ID, Valid: true}, threshold

		go subs.UnsubscribeChannelSalesUpdates()

		go func() {
			subscribedChannels := config.ActiveAllSales[threshold]
			for index, c := range subscribedChannels {
				if c == channel.ID {
					subscribedChannels = slices.Delete(subscribedChannels, index, index+1)
					config.ActiveAllSales[threshold] = subscribedChannels
					break
				}
			}
		}()

		message = fmt.Sprintf("Deactivate AllSales Subscription for Threshold : %f ETH  on Channel: %s", threshold, channel.Name)
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
