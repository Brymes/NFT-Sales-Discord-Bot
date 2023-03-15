package handlers

import (
	"DIA-NFT-Sales-Bot/config"
	"DIA-NFT-Sales-Bot/models"
	"DIA-NFT-Sales-Bot/utils"
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func TrackFloorHandler(discordSession *discordgo.Session, interaction *discordgo.InteractionCreate) {
	optionsMap := ParseCommandOptions(interaction)

	address, blockchain := optionsMap["address"].StringValue(), optionsMap["blockchain"].StringValue()

	//Respond Channel is being Setup
	message := fmt.Sprintf("Track Floor Price for Collection: %s ", utils.CreateHyperLink(address, utils.GetScanLink("address", address, blockchain)))
	err := discordSession.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})

	if err != nil {
		panic(err)
	}

	config.FloorPriceTrackerAddress = ""
	models.Subscriptions{
		Command: "track_floor_price",
		Address: sql.NullString{
			String: address,
			Valid:  true,
		},
		ChannelID: sql.NullString{
			String: interaction.GuildID,
			Valid:  true,
		},
		Blockchain: blockchain,
		Active:     true,
	}.SaveTracker()

	_ = discordSession.GuildMemberNickname(interaction.GuildID, "@me", "DIA Sales Tracker")
	config.FloorPriceTrackerAddress, config.FloorPriceTrackerChain, config.FloorPriceTrackerGuild = address, blockchain, interaction.GuildID

}
