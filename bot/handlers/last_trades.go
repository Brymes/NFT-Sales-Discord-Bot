package handlers

import (
	"DIA-NFT-Sales-Bot/config"
	"DIA-NFT-Sales-Bot/services"
	"DIA-NFT-Sales-Bot/utils"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"math"
	"time"
)

func LastTradesHandler(discordSession *discordgo.Session, interaction *discordgo.InteractionCreate) {

	channel := interaction.ChannelID

	payload, found := config.ActiveSalesInfoBot[channel]

	if !found {
		_, err := discordSession.ChannelMessageSend(channel, "This Channel is not registered. Kindly use /set_up_info_bot and select this channel.")
		if err != nil {
			panic(err)
		}
		return
	}

	//Respond Channel is being Setup
	message := fmt.Sprintf("Get Last Trades for Collection: %s ", utils.CreateHyperLink(payload["address"], utils.GetScanLink("address", payload["address"], payload["blockchain"])))
	err := discordSession.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})

	if err != nil {
		panic(err)
	}

	responses := services.LastTradesAPI(payload["address"], payload["blockchain"])
	for index, response := range responses[:10] {
		embedMsg := createLastTradesMessage(response, payload["blockchain"], payload["address"], index)
		_, err = config.DiscordBot.ChannelMessageSendEmbed(interaction.ChannelID, embedMsg)
		if err != nil {
			panic(err)
		}
	}

}

func createLastTradesMessage(payload services.LastTradesAPIResponse, blockchain, address string, count int) *discordgo.MessageEmbed {
	scanLink := utils.GetScanLink("address", address, blockchain)
	price := fmt.Sprintf("%v %s", math.Round(payload.Price*100)/100, payload.Currency.Name)
	title := fmt.Sprintf("Recent Trades of %s: No. %d ", payload.Name, count+1)

	embed := &discordgo.MessageEmbed{
		Color:       0x5f3267,
		Title:       title,
		Description: "NFT Discord Bot Last Trades",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Collection Name",
				Value:  payload.Name,
				Inline: false,
			}, {
				Name:   "Collection Address",
				Value:  utils.CreateHyperLink(address, scanLink),
				Inline: true,
			}, {
				Name:   "Price",
				Value:  price,
				Inline: true,
			}, {
				Name:   "Token ID",
				Value:  payload.NFTid,
				Inline: true,
			}, {
				Name:   "Seller Address",
				Value:  utils.CreateHyperLink(payload.FromAddress, utils.GetScanLink("address", payload.FromAddress, blockchain)),
				Inline: true,
			}, {
				Name:   "Buyer Address",
				Value:  utils.CreateHyperLink(payload.ToAddress, utils.GetScanLink("address", payload.ToAddress, blockchain)),
				Inline: true,
			}, {
				Name:   "MarketPlace",
				Value:  payload.Exchange,
				Inline: true,
			}, {
				Name:   "Time of Sale",
				Value:  payload.Timestamp.Format(time.RFC3339),
				Inline: true,
			}, {
				Name:   "TxHash",
				Value:  utils.CreateHyperLink(payload.TxHash, utils.GetScanLink("transaction", payload.TxHash, blockchain)),
				Inline: true,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
	}

	return embed
}
