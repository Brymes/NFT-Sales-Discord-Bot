package handlers

import (
	log "DIA-NFT-Sales-Bot/debug"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

func ChangeBotAvatarHandler(discordSession *discordgo.Session, interaction *discordgo.InteractionCreate) {
	optionsMap := ParseCommandOptions(interaction)

	//Respond Channel is being Setup
	message := "Currently processing Avatar change request"
	err := discordSession.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})

	if err != nil {
		panic(err)
	}

	imageID := optionsMap["image"].Value.(string)
	imageUrl := interaction.ApplicationCommandData().Resolved.Attachments[imageID].URL

	//get the file contents

	resp, err := http.DefaultClient.Get(imageUrl)

	if err != nil {
		log.Log.Println(errors.New("could not get response from code explain bot"))
		err = discordSession.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Could not get response",
			},
		})

		if err != nil {
			panic(err)
		}

	}
	defer func() {
		_ = resp.Body.Close()
	}()

	img, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Log.Println("Error reading the response, ", err)
		return
	}

	contentType := http.DetectContentType(img)
	base64img := base64.StdEncoding.EncodeToString(img)

	//Follow Up has been Set up
	avatar := fmt.Sprintf("data:%s;base64,%s", contentType, base64img)
	_, err = discordSession.UserUpdate("", avatar)
	if err != nil {
		log.Log.Println(err)
	}

	SendChannelSetupFollowUp("Avatar Successfully changed", discordSession, interaction)

}
