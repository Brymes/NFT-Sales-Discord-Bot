package handlers

import (
	"DIA-NFT-Sales-Bot/models"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/olekukonko/tablewriter"
)

// Note: could remove SetBorder and SetRowLine if char count needs to be reduced
func formatSubscriptionsText(subsArray [][]string) string {
	str := &strings.Builder{}
	table := tablewriter.NewWriter(str)
	table.SetHeader([]string{"Command", "Channel", "Threshold", "Address"})
	table.AppendBulk(subsArray)
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetHeaderLine(true)
	table.Render()
	return "```" + str.String() + "```"
}

func SubscriptionsHandler(discordSession *discordgo.Session, interaction *discordgo.InteractionCreate) {

	err := discordSession.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Loading All Subscriptions",
			Title:   "All Subscriptions request",
		},
	})

	if err != nil {
		panic(err)
	}
	var subsArray [][]string

	subscriptions := models.Subscriptions{}.LoadAllSubscriptions()

	for _, subscription := range subscriptions {
		var row []string

		switch subscription.Command {
		case "sales":
			row = []string{subscription.Command, subscription.ChannelID.String, "", subscription.Address.String}
		case "all_sales":
			row = []string{subscription.Command, subscription.ChannelID.String, subscription.Threshold.String, ""}
		}

		subsArray = append(subsArray, row)
	}

	msg := formatSubscriptionsText(subsArray)
	SendChannelSetupFollowUp(msg, discordSession, interaction)
}
