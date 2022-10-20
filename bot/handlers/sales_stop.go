package handlers

import (
	"DIA-NFT-Sales-Bot/config"
	"DIA-NFT-Sales-Bot/models"
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"strings"
)

func SalesStopHandler(discordSession *discordgo.Session, interaction *discordgo.InteractionCreate) {
	var (
		message             string
		optionsMap          = ParseCommandOptions(interaction)
		subs                = models.Subscriptions{Command: "sales"}
		address, addrExists = optionsMap["address"]
		channel, chanExists = optionsMap["channel"]
		all, blockchain     = optionsMap["all"].BoolValue(), optionsMap["blockchain"].StringValue()
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
			modifiedAddress := strings.ToUpper(address.StringValue())

			subs.ChannelID, subs.Address, subs.Blockchain = sql.NullString{String: channelID, Valid: true}, sql.NullString{String: modifiedAddress, Valid: true}, blockchain
			go subs.UnsubscribeChannelSalesUpdates()

			go func() {
				subscribedChannels := config.ActiveSales[modifiedAddress][blockchain]
				for index, c := range subscribedChannels {
					if c == channelID {
						subscribedChannels = slices.Delete(subscribedChannels, index, index+1)
						config.ActiveSales[modifiedAddress][blockchain] = subscribedChannels
						break
					}
				}
			}()

			message = fmt.Sprintf("Deactivated Sales Subscription for Contract Address : %s  on Channel: %s", modifiedAddress, channel.Name)
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
