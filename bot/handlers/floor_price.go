package handlers

import (
	"DIA-NFT-Sales-Bot/config"
	"DIA-NFT-Sales-Bot/services"
	"DIA-NFT-Sales-Bot/utils"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
	"time"
)

func FloorPriceHandler(discordSession *discordgo.Session, interaction *discordgo.InteractionCreate) {

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
	message := fmt.Sprintf("Get FloorPrice for Collection: %s ", utils.CreateHyperLink(payload["address"], utils.GetScanLink("address", payload["address"], payload["blockchain"])))
	err := discordSession.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})

	if err != nil {
		panic(err)
	}

	response := services.FloorAPI(payload["address"], payload["blockchain"])
	embedMsg := createFloorPriceMessage(response, payload["address"], payload["blockchain"])

	_, err = config.DiscordBot.ChannelMessageSendEmbed(interaction.ChannelID, embedMsg)
	if err != nil {
		panic(err)
	}
}

func createFloorPriceMessage(payload services.FloorPriceResponse, address, blockchain string) *discordgo.MessageEmbed {
	scanLink := utils.GetScanLink("address", address, blockchain)

	embed := &discordgo.MessageEmbed{
		Color:       0x5f3267,
		Title:       "Floor Price",
		Description: "NFT Discord Bot Floor Price Response",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Collection Address",
				Value:  utils.CreateHyperLink(address, scanLink),
				Inline: true,
			}, {
				Name:   "Collection FloorPrice",
				Value:  fmt.Sprintf("%f %s", payload.FloorPrice, currencies[strings.ToLower(blockchain)]),
				Inline: true,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
	}

	return embed
}
