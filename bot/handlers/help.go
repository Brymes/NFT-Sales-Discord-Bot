package handlers

import (
	"github.com/bwmarrin/discordgo"
	"strings"

	"github.com/olekukonko/tablewriter"
)

var (
	HelpText  = formatHelpText(helpArray[:7])
	HelpText2 = formatHelpText(helpArray[8:])
	helpArray = [][]string{
		{"/help", "Display help message"},
		{"/subscriptions", "Returns a list of commands which the server has enabled"},
		{"/sales", "Get all sales for provided NFT collection to selected discord channel"},
		{"/sales_stop", "Stops Bot from pushing sales update for NFT or stop all sales bots"},
		{"/floor", "Return floor price of provided NFT collection address"},
		{"/all_sales", "Return all sales above the provided threshold"},
		{"/all_sales_stop", "Stop bot for all sales above the set threshold and collection"},
		{"/stop_all", "Stops all bots from operating in the selected channel"},
		{"/set_up_info_bot", "Enable channel to process utility commands to track NFT sales"},
		{"/volume", "Get volume for NFT  Collection set via set_up_info_bot"},
		{"/floor_price", "Get floor price for NFT Collection set via set_up_info_bot"},
		{"/last_trades", "Get recent trades for NFT collection set via set_up_info_bot"},
		{"/stop_subscription", "Returns modal to select commands to stop"},
	}
)

// Note: could remove SetBorder and SetRowLine if char count needs to be reduced
func formatHelpText(helpArray [][]string) string {
	str := &strings.Builder{}
	table := tablewriter.NewWriter(str)
	table.SetHeader([]string{"Help", "Message"})
	table.AppendBulk(helpArray)
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetHeaderLine(true)
	table.Render()
	return "```" + str.String() + "```"
}

func HelpHandler(discordSession *discordgo.Session, interaction *discordgo.InteractionCreate) {
	err := discordSession.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: HelpText,
			Title:   "Help/Commands Message",
		},
	})

	if err != nil {
		panic(err)
	}
	SendChannelSetupFollowUp(HelpText2, discordSession, interaction)
}

func SendHelpText(discordSession *discordgo.Session, message *discordgo.MessageCreate) {
	_, err := discordSession.ChannelMessageSend(message.ChannelID, HelpText)
	_, err = discordSession.ChannelMessageSend(message.ChannelID, HelpText2)
	if err != nil {
		panic(err)
	}
}
