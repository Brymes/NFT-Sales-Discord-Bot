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

func StopAllHandler(discordSession *discordgo.Session, interaction *discordgo.InteractionCreate) {
	optionsMap := ParseCommandOptions(interaction)
	defer config.ActiveSalesMux.Unlock()
	defer config.ActiveAllSalesMux.Unlock()
	defer config.ActiveSalesInfoMux.Unlock()

	if optionsMap["channel"] != nil {
		channel := optionsMap["channel"].ChannelValue(discordSession)

		err := discordSession.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Stop all Bots for Channel %s", channel.Name),
			},
		})

		if err != nil {
			panic(err)
		}

		sub := models.Subscriptions{ChannelID: sql.NullString{
			String: channel.ID,
			Valid:  true,
		}}

		channelSubs := sub.LoadChannelSubscriptions()

		for _, channelSub := range channelSubs {
			switch channelSub.Command {
			case "sales":
				for chainName, chain := range config.ActiveSales[channelSub.Address.String] {
					for index, subChannel := range chain {
						if subChannel == channelSub.ChannelID.String {
							config.ActiveSalesMux.Lock()
							config.ActiveSales[channelSub.Address.String][chainName] = slices.Delete(config.ActiveSales[channelSub.Address.String][chainName], index, index+1)
							break
						}
					}
				}
			case "all_sales":
				for chainName, chain := range config.ActiveAllSales[channelSub.Threshold] {
					for index, subChannel := range chain {
						if subChannel == channelSub.ChannelID.String {
							config.ActiveAllSalesMux.Lock()
							config.ActiveAllSales[channelSub.Threshold][chainName] = slices.Delete(config.ActiveAllSales[channelSub.Threshold][chainName], index, index+1)
							break
						}
					}
				}
			}
		}
		SendChannelSetupFollowUp("Done stopping bots for Selected channel", discordSession, interaction)

		go sub.DeactivateChannelSubscriptions()
	} else {

		err := discordSession.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Stop all Bots",
			},
		})

		if err != nil {
			panic(err)
		}

		go models.Subscriptions{}.DeactivateAllSubscriptions()
		StopAllBots()

		SendChannelSetupFollowUp("Done stopping bots for all channels", discordSession, interaction)

	}

}

func StopAllBots() {
	go config.ShutDownWS()

	if !config.ActiveSalesMux.TryLock() {
		config.ActiveSalesMux.Unlock()
		config.ActiveSalesMux.Lock()
	}
	if !config.ActiveAllSalesMux.TryLock() {
		config.ActiveAllSalesMux.Unlock()
		config.ActiveAllSalesMux.Lock()
	}
	if !config.ActiveSalesInfoMux.TryLock() {
		config.ActiveSalesInfoMux.Unlock()
		config.ActiveSalesInfoMux.Lock()
	}

	// Delete Global variables
	config.ActiveAllSalesKeys = nil
	maps.Clear(config.ActiveAllSales)
	maps.Clear(config.ActiveSales)
	maps.Clear(config.ActiveSalesInfoBot)
}
