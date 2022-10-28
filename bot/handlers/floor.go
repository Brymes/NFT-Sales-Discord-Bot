package handlers

import (
	"DIA-NFT-Sales-Bot/config"
	"DIA-NFT-Sales-Bot/services"
	"DIA-NFT-Sales-Bot/utils"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func FloorHandler(discordSession *discordgo.Session, interaction *discordgo.InteractionCreate) {
	optionsMap := ParseCommandOptions(interaction)

	address, blockchain := optionsMap["address"].StringValue(), optionsMap["blockchain"].StringValue()

	//Respond Channel is being Setup
	message := fmt.Sprintf("Get Floor Price for Collection: %s ", utils.CreateHyperLink(address, utils.GetScanLink("address", address, blockchain)))
	err := discordSession.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})

	if err != nil {
		panic(err)
	}

	response := services.FloorPriceAPI(address)
	embedMsg := createFloorMessage(response, blockchain)

	_, err = config.DiscordBot.ChannelMessageSendEmbed(interaction.ChannelID, embedMsg)
	if err != nil {
		panic(err)
	}
}

func createFloorMessage(payload services.Floor, blockchain string) *discordgo.MessageEmbed {
	etherScanLink := utils.GetScanLink("address", payload.Volume.Address, blockchain)

	title := fmt.Sprintf("Floor Price of %s is  %f", payload.Volume.Collection, payload.FloorPrice.FloorPrice)

	embed := &discordgo.MessageEmbed{
		Color:       0x5f3267,
		Title:       title,
		Description: "NFT Discord Bot Floor Price Response",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Collection Name",
				Value:  payload.Volume.Collection,
				Inline: false,
			}, {
				Name:   "Collection Address",
				Value:  utils.CreateHyperLink(payload.Volume.Address, etherScanLink),
				Inline: true,
			}, {
				Name:   "Floor Price",
				Value:  fmt.Sprintf("%f %s", payload.FloorPrice.FloorPrice, currencies[strings.ToLower(blockchain)]),
				Inline: false,
			}, {
				Name:   "Moving Average Floor Price",
				Value:  fmt.Sprintf("%f %s", payload.MA.MovingAverageFloorPrice, currencies[strings.ToLower(blockchain)]),
				Inline: false,
			}, {
				Name:   "24h Volume",
				Value:  fmt.Sprintf("%f %s", payload.Volume.Volume, currencies[strings.ToLower(blockchain)]),
				Inline: false,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
		Footer:    &config.MessageFooter,
	}

	return embed
}
