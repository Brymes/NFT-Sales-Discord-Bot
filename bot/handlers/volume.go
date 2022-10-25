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

func VolumeHandler(discordSession *discordgo.Session, interaction *discordgo.InteractionCreate) {

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
	message := fmt.Sprintf("Get Volume for Collection: %s ", utils.GetScanLink("address", payload["address"], payload["blockchain"]))
	err := discordSession.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})

	if err != nil {
		panic(err)
	}

	response := services.VolumeAPI(payload["address"], payload["blockchain"])
	embedMsg := createVolumeMessage(response, payload["blockchain"])

	_, err = config.DiscordBot.ChannelMessageSendEmbed(interaction.ChannelID, embedMsg)
	if err != nil {
		panic(err)
	}
}

func createVolumeMessage(payload services.VolumeAPIResponse, blockchain string) *discordgo.MessageEmbed {
	scanLink := utils.GetScanLink("address", payload.Address, blockchain)

	title := fmt.Sprintf("Volume of %s is  %f %s", payload.Collection, payload.Volume, currencies[strings.ToLower(blockchain)])

	embed := &discordgo.MessageEmbed{
		Color:       0x5f3267,
		Title:       title,
		Description: "NFT Discord Bot Floor Price Response",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Collection Name",
				Value:  payload.Collection,
				Inline: false,
			}, {
				Name:   "Collection Address",
				Value:  utils.CreateHyperLink(payload.Address, scanLink),
				Inline: true,
			}, {
				Name:   "Collection 24h Volume",
				Value:  fmt.Sprintf("%f", payload.Volume),
				Inline: true,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
	}

	return embed
}
