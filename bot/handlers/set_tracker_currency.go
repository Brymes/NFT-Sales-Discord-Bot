package handlers

import (
	"DIA-NFT-Sales-Bot/config"
	"DIA-NFT-Sales-Bot/models"
	"github.com/bwmarrin/discordgo"
)

func SetCurrencyHandler(discordSession *discordgo.Session, interaction *discordgo.InteractionCreate) {
	optionsMap := ParseCommandOptions(interaction)

	currency := optionsMap["currency"].StringValue()

	_ = discordSession.GuildMemberNickname(interaction.GuildID, "@me", "DIA Sales Tracker")
	config.TrackerCurrency = currency
	cm := models.ConfigModel{TrackerCurrency: currency}
	cm.SaveConfig()
}
