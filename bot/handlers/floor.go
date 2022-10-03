package handlers

import (
	"DIA-NFT-Sales-Bot/config"
	"DIA-NFT-Sales-Bot/services"
	"DIA-NFT-Sales-Bot/utils"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

func FloorHandler(discordSession *discordgo.Session, interaction *discordgo.InteractionCreate) {
	optionsMap := ParseCommandOptions(interaction)

	address := optionsMap["address"].StringValue()

	//Respond Channel is being Setup
	message := fmt.Sprintf("Get Floor Price for Collection: %s ", utils.GetEtherScanLink("address", address))
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
	embedMsg := createFloorMessage(response)

	_, err = config.DiscordBot.ChannelMessageSendEmbed(interaction.ChannelID, embedMsg)
	if err != nil {
		panic(err)
	}
}

func createFloorMessage(payload services.Floor) *discordgo.MessageEmbed {
	etherScanLink := utils.GetEtherScanLink("address", payload.Volume.Address)

	title := fmt.Sprintf("Floor Price of %s is  %f ETH", payload.Volume.Collection, payload.FloorPrice.FloorPrice)

	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0x5f3267,
		Title:       title,
		Description: "NFT Discord Bot Floor Price Response",
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "Collection Name",
				Value:  payload.Volume.Collection,
				Inline: false,
			}, &discordgo.MessageEmbedField{
				Name:   "Collection Address",
				Value:  etherScanLink,
				Inline: true,
			}, &discordgo.MessageEmbedField{
				Name:   "Floor Price",
				Value:  fmt.Sprintf("%f ETH", payload.FloorPrice.FloorPrice),
				Inline: false,
			}, &discordgo.MessageEmbedField{
				Name:   "Moving Average Floor Price",
				Value:  fmt.Sprintf("%f ETH", payload.FloorPrice.FloorPrice),
				Inline: false,
			}, &discordgo.MessageEmbedField{
				Name:   "24h Volume",
				Value:  fmt.Sprintf("%f ETH", payload.Volume.Volume),
				Inline: false,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
		//Timestamp: eventResponse.Timestamp.Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
	}

	return embed
}
