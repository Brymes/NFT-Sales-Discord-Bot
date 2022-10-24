package handlers

import (
	"DIA-NFT-Sales-Bot/config"
	"DIA-NFT-Sales-Bot/models"
	"DIA-NFT-Sales-Bot/services"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strconv"
)

func StopSubscriptionsHandler(discordSession *discordgo.Session, interaction *discordgo.InteractionCreate) {

	//Respond Channel is being Setup
	message := "Loading All Subscriptions"
	err := discordSession.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})

	if err != nil {
		panic(err)
	}

	subscriptions := parseSubscriptionsFromDB(discordSession)
	minimumValues := 1
	response := discordgo.WebhookParams{
		Content: "Kindly Select from the list below bots you would like to kill",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						CustomID:    "commands_to_stop",
						Placeholder: "Stop Bots",
						MinValues:   &minimumValues,
						MaxValues:   10,
						Options:     subscriptions,
					},
				},
			},
		},
	}

	SendComplexFollowUp(response, discordSession, interaction)

}

func parseSubscriptionsFromDB(discordSession *discordgo.Session) (subs []discordgo.SelectMenuOption) {
	allSubscriptions := models.Subscriptions{}.LoadAllSubscriptions()
	subs = []discordgo.SelectMenuOption{}

	for _, subscription := range allSubscriptions {

		channel, err := discordSession.State.Channel(subscription.ChannelID.String)

		if err != nil {
			channel, _ = discordSession.Channel(subscription.ChannelID.String)
		}

		label, description := fmt.Sprintf("Comand %s on channel: %s", subscription.Command, channel.Name), ""

		if subscription.Address.Valid {
			description = fmt.Sprintf("For Address: %s on Blockchain: %s", subscription.Address.String, subscription.Blockchain)
		} else {
			description = fmt.Sprintf("For Threshold: %f on Blockchain: %s", subscription.Threshold, subscription.Blockchain)
		}

		subs = append(subs, discordgo.SelectMenuOption{
			Label:       label,
			Value:       fmt.Sprintf("%d", subscription.ID),
			Description: description,
			Emoji:       discordgo.ComponentEmoji{Name: "❌"},
			Default:     false,
		})
	}

	return subs
}

func StopSubscriptions(discordSession *discordgo.Session, interaction *discordgo.InteractionCreate) {
	data := interaction.MessageComponentData()

	err := discordSession.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Deleting Selected Subscriptions",
		},
	})

	if err != nil {
		panic(err)
	}

	idArray := make([]int, len(data.Values))
	for index, dbIndex := range data.Values {
		num, _ := strconv.Atoi(dbIndex)
		idArray[index] = num
	}

	go models.Subscriptions{}.DeactivateSubscriptions(idArray)

	startWS := models.LoadCurrentSubscriptions()
	if !config.ActiveNftEventWS && startWS {
		services.StartEventWS()
	}

	message := fmt.Sprintf("Successfully Unsubscribed from bots on %d channels", len(idArray))
	SendChannelSetupFollowUp(message, discordSession, interaction)
}
