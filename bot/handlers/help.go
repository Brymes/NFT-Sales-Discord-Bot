package handlers

import (
	"github.com/bwmarrin/discordgo"
	"strings"

	"github.com/olekukonko/tablewriter"
)

var (
	HelpText  = formatHelpText()
	helpArray = [][]string{
		{"/help", "Returns All Commands and their corresponding Descriptions"},
		{"/subscriptions", "Returns a list of commands which the server has enabled"},
		{"/sales", "Allow to set-up a channel that feeds all NFT collection sales matching the supplied contract address from the DIA NFT Event WebSocket to a selected discord channel"},
		{"/sales_stop", "Stops Bot from pushing sales update from a contract address or stop all bots"},
		{"/floor", "Return floor price of the provided NFT collection contract address"},
		{"/all_sales", "Return all sales above the provided threshold"},
		{"/all_sales_stop", "Stop bot for all sales above the predetermined threshold and contract address"},
		{"/stop_all", "Stops all bots from operating in the selected channel or stop all bots if channel is not provided"},
	}
)

// Note: could remove SetBorder and SetRowLine if char count needs to be reduced
func formatHelpText() string {
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
}

func SendHelpText(discordSession *discordgo.Session, message *discordgo.MessageCreate) {
	_, err := discordSession.ChannelMessageSend(message.ChannelID, HelpText)
	if err != nil {
		panic(err)
	}
}
