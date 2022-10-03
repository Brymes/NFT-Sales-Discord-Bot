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

func SalesStopHandler(discordSession *discordgo.Session, interaction *discordgo.InteractionCreate) {
	var (
		message             string
		optionsMap          = ParseCommandOptions(interaction)
		subs                = models.Subscriptions{Command: "sales"}
		address, addrExists = optionsMap["address"]
		channel, chanExists = optionsMap["channel"]
		all                 = optionsMap["all"].BoolValue()
	)
	config.ActiveSalesMux.Lock()

	//Respond Channel is being Setup
	if all {
		message = "Deactivate Sales Subscription for all Channels and Contract Addresses"

		go maps.Clear(config.ActiveSales)
		go subs.UnsubscribeSalesUpdates()

	} else {
		if !addrExists || !chanExists {
			message = "Invalid Channel or Address supplied"
		} else {
			channelID := channel.ChannelValue(discordSession).ID

			subs.ChannelID, subs.Address = sql.NullString{String: channelID, Valid: true}, sql.NullString{String: address.StringValue(), Valid: true}
			go subs.UnsubscribeChannelSalesUpdates()

			go func() {
				subscribedChannels := config.ActiveSales[address.StringValue()]
				for index, c := range subscribedChannels {
					if c == channelID {
						subscribedChannels = slices.Delete(subscribedChannels, index, index+1)
						config.ActiveSales[address.StringValue()] = subscribedChannels
						break
					}
				}
			}()

			message = fmt.Sprintf("Deactivated Sales Subscription for Contract Address : %s  on Channel: %s", address.StringValue(), channel.Name)
		}
	}
	defer config.ActiveSalesMux.Unlock()

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
