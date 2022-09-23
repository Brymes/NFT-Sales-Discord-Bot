package services

import (
	"DIA-NFT-Sales-Bot/config"
	"DIA-NFT-Sales-Bot/utils"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
	"time"
)

func StartEventWS() {
	fileName := fmt.Sprintf("WSLogger-%s", strings.Split(time.Now().String(), " ")[0])
	WSLogger := config.CreateServiceLogger(fileName)
	reqBuffer, logger := config.InitRequestLogger("NftEvents")

	defer WSLogger.Printf("\n%v\n", reqBuffer)
	defer utils.HandlePanic(logger, "")

	// Use Cancel Func to kill this
	go ConnectToService(logger)
	return
}

func SalesController(event *NFTEvent) {
	go HandleSales(event)
	go HandleAllSales(event)
}

func HandleSales(event *NFTEvent) {
	config.ActiveSalesMux.RLock()
	channels, match := config.ActiveSales[event.Response.NFT.NFTClass.Address]

	if !match {
		return
	} else {
		for _, channel := range channels {
			go SendSalesMessage(event, channel)
		}
	}
	config.ActiveSalesMux.RUnlock()
}

func HandleAllSales(event *NFTEvent) {
	priceInEth := utils.ConvertDecimalsToEth(event.Response.Price, event.Response.Currency.Decimals)

	config.ActiveAllSalesMux.RLock()

	for _, elem := range config.ActiveAllSalesKeys {
		if priceInEth > elem {
			for _, channel := range config.ActiveAllSales[elem] {
				go SendSalesMessage(event, channel)
			}
		}
	}

	config.ActiveAllSalesMux.RUnlock()

}

func SendSalesMessage(event *NFTEvent, channelID string) {
	eventResponse := event.Response
	priceInEth := fmt.Sprintf("%v", utils.ConvertDecimalsToEth(eventResponse.Price, eventResponse.Currency.Decimals))
	marketPlaceLink := utils.GetMarketPlaceLink(eventResponse.Exchange, eventResponse.NFT.NFTClass.Address, eventResponse.NFT.TokenID)
	title := fmt.Sprintf("NFT Sale @ %s ETH on %s", priceInEth, eventResponse.Exchange)

	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0x5f3267,
		Title:       title,
		Description: "This is a discordgo embed",
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "Collection Name",
				Value:  marketPlaceLink,
				Inline: false,
			}, &discordgo.MessageEmbedField{
				Name:   "Seller Address",
				Value:  utils.GetEtherScanLink("address", eventResponse.FromAddress),
				Inline: true,
			}, &discordgo.MessageEmbedField{
				Name:   "Buyer Address",
				Value:  utils.GetEtherScanLink("address", eventResponse.ToAddress),
				Inline: true,
			}, &discordgo.MessageEmbedField{
				Name:   "Price In Eth",
				Value:  priceInEth,
				Inline: true,
			}, &discordgo.MessageEmbedField{
				Name:   "MarketPlace",
				Value:  eventResponse.Exchange,
				Inline: true,
			}, &discordgo.MessageEmbedField{
				Name:   "Time of Sale",
				Value:  eventResponse.Timestamp.Format(time.RFC3339),
				Inline: true,
			}, &discordgo.MessageEmbedField{
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
