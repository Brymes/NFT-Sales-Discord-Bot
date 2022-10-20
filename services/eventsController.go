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
	chains, match := config.ActiveSales[strings.ToUpper(event.Response.NFT.NFTClass.Address)]
	config.ActiveSalesMux.Unlock()

	if !match {
		return
	} else {
		channels := chains[event.Response.NFT.NFTClass.Blockchain]
		for _, channel := range channels {
			go SendSalesMessage(event, channel)
		}
	}
}

func HandleAllSales(event NFTEvent) {
	// This Handle panic is useful for if all bots are stopped and arrays/maps have been emptied
	defer utils.HandlePanic(config.DiscordBot, "Error in All Sales Handler")
	price := utils.ConvertDecimalsToCurrency(event.Response.Price, event.Response.Currency.Decimals)

	config.ActiveAllSalesMux.Lock()
	for _, elem := range config.ActiveAllSalesKeys {
		if price > elem {
			chains := config.ActiveAllSales[elem]
			for _, channel := range chains[event.Response.NFT.NFTClass.Blockchain] {
				go SendSalesMessage(event, channel)
			}
		}
	}

	defer config.ActiveAllSalesMux.Unlock()

}

func SendSalesMessage(event NFTEvent, channelID string) {
	defer utils.HandlePanic(config.DiscordBot, "Error Sending sales event")
	var price, txHash, buyersAddress, sellerAddress string

	eventResponse := event.Response
	blockchain := eventResponse.NFT.NFTClass.Blockchain

	switch blockchain {
	case "Ethereum":
		price, txHash, buyersAddress, sellerAddress = parseEthereumSalesMessage(event)
	case "Solana":
		price, txHash, buyersAddress, sellerAddress = parseSolanaSalesMessage(event)
	case "Astar":
		price, txHash, buyersAddress, sellerAddress = parseAstarSalesMessage(event)
	}

	marketPlaceLink := utils.GetMarketPlaceLink(eventResponse.Exchange, eventResponse.NFT.NFTClass.Address, eventResponse.NFT.TokenID)
	title := fmt.Sprintf("NFT Sale @ %s on %s", price, eventResponse.Exchange)

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
				Value:  sellerAddress,
				Inline: true,
			}, {
				Name:   "Buyer Address",
				Value:  buyersAddress,
				Inline: true,
			}, {
				Name:   "Price",
				Value:  price,
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
				Value:  txHash,
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

func parseEthereumSalesMessage(event NFTEvent) (price, txHash, buyersAddress, sellerAddress string) {
	eventResponse := event.Response
	sellerAddress = utils.GetScanLink("address", eventResponse.FromAddress, "ethereum")
	buyersAddress = utils.GetScanLink("address", eventResponse.ToAddress, "ethereum")
	price = fmt.Sprintf("%v ETH", utils.ConvertDecimalsToCurrency(eventResponse.Price, eventResponse.Currency.Decimals))
	txHash = utils.GetScanLink("transaction", eventResponse.TxHash, "ethereum")
	return
}

func parseSolanaSalesMessage(event NFTEvent) (price, txHash, buyersAddress, sellerAddress string) {
	eventResponse := event.Response
	sellerAddress = utils.GetScanLink("address", eventResponse.FromAddress, "solana")
	buyersAddress = utils.GetScanLink("address", eventResponse.ToAddress, "solana")
	price = fmt.Sprintf("%v SOL", utils.ConvertDecimalsToCurrency(eventResponse.Price, eventResponse.Currency.Decimals))
	txHash = utils.GetScanLink("transaction", eventResponse.TxHash, "solana")
	return
}

func parseAstarSalesMessage(event NFTEvent) (price, txHash, buyersAddress, sellerAddress string) {
	eventResponse := event.Response
	sellerAddress = utils.GetScanLink("address", eventResponse.FromAddress, "astar")
	buyersAddress = utils.GetScanLink("address", eventResponse.ToAddress, "astar")
	price = fmt.Sprintf("%v ASTR", utils.ConvertDecimalsToCurrency(eventResponse.Price, eventResponse.Currency.Decimals))
	txHash = utils.GetScanLink("transaction", eventResponse.TxHash, "astar")
	return
}
