package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var (
	ContractAddressMinLength = 42
	RegisteredCommands       []*discordgo.ApplicationCommand
	TextChannelType          = []discordgo.ChannelType{
		discordgo.ChannelTypeGuildText,
		discordgo.ChannelTypeGroupDM,
		discordgo.ChannelTypeDM,
		discordgo.ChannelTypeGuildNews,
		discordgo.ChannelTypeGuildNewsThread,
		discordgo.ChannelTypeGuildCategory,
		discordgo.ChannelTypeGuildPrivateThread,
		discordgo.ChannelTypeGuildPublicThread,
		discordgo.ChannelTypeGuildStore,
	}
	BlockChainChoices = []*discordgo.ApplicationCommandOptionChoice{
		{
			Name:  "Astar",
			Value: "Astar",
		}, {
			Name:  "Ethereum",
			Value: "Ethereum",
		},
		// {
		// 	Name:  "Solana",
		// 	Value: "Solana",
		// },
	}
)

func RegisterCommands(discordSession *discordgo.Session) {
	var commands = []*discordgo.ApplicationCommand{
		{
			Name:        "help",
			Description: "Returns All Commands and Descriptions",
		},
		{
			Name:        "subscriptions",
			Description: "Returns a list of commands which the server has enabled",
		},
		{
			Name:        "stop_subscription",
			Description: "Select which commands/bots kill",
		}, {
			Name:        "volume",
			Description: "Returns volume information for the predetermined NFT Collection on Info bot",
		}, {
			Name:        "floor_price",
			Description: "Returns floor price information for the predetermined NFT Collection on Info bot",
		}, {
			Name:        "last_trades",
			Description: "Returns last 10 trades information for the predetermined NFT Collection on Info bot",
		},
		{
			Name:        "set_up_info_bot",
			Description: "Set up bot for volume, floor_price and last_trades information with single command for an collection",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "collection_address",
					Description: "Contract Address to filter transactions from",
					Required:    true,
					MinLength:   &ContractAddressMinLength,
				},
				{
					Type:         discordgo.ApplicationCommandOptionChannel,
					Name:         "channel",
					Description:  "Channel to push information of matching transactions to.",
					ChannelTypes: TextChannelType,
					Required:     true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "blockchain",
					Description: "Select from Astar or Ethereum",
					Choices:     BlockChainChoices,
					Required:    true,
				},
			},
		}, {
			Name:        "sales",
			Description: "Set up bot that feeds all sales for a selected NFT Collection to the channel of your choice",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "collection_address",
					Description: "Contract Address to filter transactions from",
					Required:    true,
					MinLength:   &ContractAddressMinLength,
				},
				{
					Type:         discordgo.ApplicationCommandOptionChannel,
					Name:         "channel",
					Description:  "Channel to push information of matching transactions to.",
					ChannelTypes: TextChannelType,
					Required:     true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "blockchain",
					Description: "Select from Astar or Ethereum",
					Choices:     BlockChainChoices,
					Required:    true,
				},
			},
		},
		// 		{
		// 			Name:        "sales_stop",
		// 			Description: "Stops Bot from pushing sales update from a contract address or stop all bots",
		// 			Options: []*discordgo.ApplicationCommandOption{
		// 				{
		// 					Type:        discordgo.ApplicationCommandOptionBoolean,
		// 					Name:        "all",
		// 					Description: "Select True to stop all sales bots",
		// 					Required:    true,
		// 				},
		// 				{
		// 					Type:        discordgo.ApplicationCommandOptionString,
		// 					Name:        "address",
		// 					Description: "Contract Address to filter transactions from",
		// 					Required:    false,
		// 					MinLength:   &ContractAddressMinLength,
		// 				},
		// 				{
		// 					Type:         discordgo.ApplicationCommandOptionChannel,
		// 					Name:         "channel",
		// 					Description:  "Channel to stop updating",
		// 					ChannelTypes: TextChannelType,
		// 					Required:     false,
		// 				},
		// 			},
		// 		},
		// 		{
		// 			Name:        "floor",
		// 			Description: "Get floor price for any NFT collection",
		// 			Options: []*discordgo.ApplicationCommandOption{

		// 				{
		// 					Type:        discordgo.ApplicationCommandOptionString,
		// 					Name:        "collection_address",
		// 					Description: "Contract Address to retrieve floor price",
		// 					Required:    true,
		// 					MinLength:   &ContractAddressMinLength,
		// 				},
		// 				{
		// 					Type:        discordgo.ApplicationCommandOptionString,
		// 					Name:        "blockchain",
		// 					Description: "Select from Astar or Ethereum",
		// 					Choices:     BlockChainChoices,
		// 					Required:    true,
		// 				},
		// 			},
		// 		},
		{
			Name:        "all_sales",
			Description: "Set up bot that feeds all NFT sales above the predetermined threshold",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionNumber,
					Name:        "threshold",
					Description: "Threshold in ETH up to 2 decimals e.g 4.55",
					Required:    true,
				},
				{
					Type:         discordgo.ApplicationCommandOptionChannel,
					Name:         "channel",
					Description:  "Channel to push information of matching transactions to.",
					ChannelTypes: TextChannelType,
					Required:     true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "blockchain",
					Description: "Select from Astar or Ethereum",
					Choices:     BlockChainChoices,
					Required:    true,
				},
			},
		},
		// 		{
		// 			Name:        "all_sales_stop",
		// 			Description: "Stop bot for all sales above the predetermined threshold and contract address",
		// 			Options: []*discordgo.ApplicationCommandOption{
		// 				{
		// 					Type:        discordgo.ApplicationCommandOptionBoolean,
		// 					Name:        "all",
		// 					Description: "Select True to stop all threshold bots",
		// 					Required:    true,
		// 				},
		// 				{
		// 					Type:        discordgo.ApplicationCommandOptionString,
		// 					Name:        "blockchain",
		// 					Description: "Select from Astar or Ethereum",
		// 					Choices:     BlockChainChoices,
		// 					Required:    true,
		// 				},
		// 				{
		// 					Type:        discordgo.ApplicationCommandOptionNumber,
		// 					Name:        "threshold",
		// 					Description: "Threshold in ETH up to 2 decimals e.g 4.55",
		// 					Required:    false,
		// 				},
		// 				{
		// 					Type:         discordgo.ApplicationCommandOptionChannel,
		// 					Name:         "channel",
		// 					Description:  "Channel to push information of matching transactions to.",
		// 					ChannelTypes: TextChannelType,
		// 					Required:     false,
		// 				},
		// 			},
		// 		},
		{
			Name:        "stop_all",
			Description: "Stops all bots from operating in the selected channel or stop all bots if channel is not provided",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionChannel,
					Name:         "channel",
					Description:  "Channel to push information of matching transactions to.",
					ChannelTypes: TextChannelType,
					Required:     false,
				},
			},
		},
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	log.Println("Added commands...", len(commands))

	for index, command := range commands {
		cmd, err := discordSession.ApplicationCommandCreate(discordSession.State.User.ID, "", command)
		if err != nil {
			log.Fatalf("Cannot create '%v' command: %v", command.Name, err)
		}
		registeredCommands[index] = cmd
	}
	RegisteredCommands = registeredCommands
}

func DeRegisterCommands(discordSession *discordgo.Session) {
	log.Println("Removing commands...")
	// // We need to fetch the commands, since deleting requires the command ID.
	// // We are doing this from the returned commands on line 375, because using
	// // this will delete all the commands, which might not be desirable, so we
	// // are deleting only the commands that we added.
	// registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
	// if err != nil {
	// 	log.Fatalf("Could not fetch registered commands: %v", err)
	// }

	for _, command := range RegisteredCommands {
		err := discordSession.ApplicationCommandDelete(discordSession.State.User.ID, "", command.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", command.Name, err)
		}
	}
}
