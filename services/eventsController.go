package services

import (
	"DIA-NFT-Sales-Bot/config"
	"DIA-NFT-Sales-Bot/utils"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
	"time"
)

func StartEventWS() {
	fileName := fmt.Sprintf("WSLogger-%s", strings.Split(time.Now().String(), " ")[0])
	WSLogger := config.CreateServiceLogger(fileName)
	reqBuffer, logger := config.InitRequestLogger("NftEvents")

	defer WSLogger.Printf("\n%v\n", reqBuffer)
	defer utils.HandlePanic(config.DiscordBot, "Error from Websocket")

	// Use Cancel Func to kill this
	log.Println("Websocket Service running")
	go ConnectToService(logger)
	return
}

func SalesController(event NFTEvent) {
	go HandleSales(event)
	go HandleAllSales(event)
}

func HandleSales(event NFTEvent) {
	// This Handle panic is useful for if all bots are stopped and arrays/maps have been emptied
	defer utils.HandlePanic(config.DiscordBot, "Error in Sales Handler")

	config.ActiveSalesMux.Lock()
	channels, match := config.ActiveSales[strings.ToUpper(event.Response.NFT.NFTClass.Address)]
	config.ActiveSalesMux.Unlock()

	if !match {
		return
	} else {
		for _, channel := range channels {
			go SendSalesMessage(event, channel)
		}
	}
}

func HandleAllSales(event NFTEvent) {
	// This Handle panic is useful for if all bots are stopped and arrays/maps have been emptied
	defer utils.HandlePanic(config.DiscordBot, "Error in All Sales Handler")
	priceInEth := utils.ConvertDecimalsToEth(event.Response.Price, event.Response.Currency.Decimals)

	config.ActiveAllSalesMux.Lock()

	for _, elem := range config.ActiveAllSalesKeys {
		if priceInEth > elem {
			for _, channel := range config.ActiveAllSales[elem] {
				go SendSalesMessage(event, channel)
			}
		}
	}
	config.ActiveAllSalesMux.Unlock()
}

func SendSalesMessage(event NFTEvent, channelID string) {
	defer utils.HandlePanic(config.DiscordBot, "Error Sending sales event")
	eventResponse := event.Response
	priceInEth := fmt.Sprintf("%v", utils.ConvertDecimalsToEth(eventResponse.Price, eventResponse.Currency.Decimals))
	marketPlaceLink := utils.GetMarketPlaceLink(eventResponse.Exchange, eventResponse.NFT.NFTClass.Address, eventResponse.NFT.TokenID)
	title := fmt.Sprintf("NFT Sale @ %s ETH on %s", priceInEth, eventResponse.Exchange)

	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0x5f3267,
		Title:       title,
		Description: "NFT Discord Bot Sales Notification",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Collection Name",
				Value:  marketPlaceLink,
				Inline: false,
			}, {
				Name:   "Seller Address",
				Value:  utils.GetEtherScanLink("address", eventResponse.FromAddress),
				Inline: true,
			}, {
				Name:   "Buyer Address",
				Value:  utils.GetEtherScanLink("address", eventResponse.ToAddress),
				Inline: true,
			}, {
				Name:   "Price In Eth",
				Value:  priceInEth,
				Inline: true,
			}, {
				Name:   "MarketPlace",
				Value:  eventResponse.Exchange,
				Inline: true,
			}, {
				Name:   "Time of Sale",
				Value:  eventResponse.Timestamp.Format(time.RFC3339),
				Inline: true,
			}, {
				Name:   "TxHash",
				Value:  utils.GetEtherScanLink("transaction", eventResponse.TxHash),
				Inline: true,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
		//Timestamp: eventResponse.Timestamp.Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
	}

	_, err := config.DiscordBot.ChannelMessageSendEmbed(channelID, embed)
	if err != nil {
		panic(err)
	}

}
