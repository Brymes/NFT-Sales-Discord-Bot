package config

import (
	"context"
	"math/big"
	"os"
	"sync"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

var (
	DBClient *gorm.DB

	DiscordBot *discordgo.Session

	// ActiveSales Maps Address to A List of channels
	ActiveSales    = map[string]map[string][]string{}
	ActiveSalesMux = &sync.Mutex{}

	// ActiveAllSales Maps Threshold to A List of Channels
	ActiveAllSales     = map[*big.Float]map[string][]string{}
	ActiveAllSalesKeys []*big.Float
	ActiveAllSalesMux  = &sync.Mutex{}

	ActiveNftEventWS     = false
	NftEventWSCancelFunc context.CancelFunc

	ActiveSalesInfoBot = map[string]map[string]string{}
	ActiveSalesInfoMux = &sync.Mutex{}

	// PanicChannelID Variable to Hold ChannelID to forward all errors
	PanicChannelID string

	FloorPriceTrackerAddress string
	FloorPriceTrackerChain   string
	FloorPriceTrackerGuild   string

	MessageFooter = discordgo.MessageEmbedFooter{
		Text:    "Powered by DIAdata.org",
		IconURL: "https://www.diadata.org",
	}
)

func ShutDownWS() {
	if ActiveNftEventWS {
		NftEventWSCancelFunc()
	}
}

func InitPanicChannel() {
	channel := os.Getenv("PANIC_CHANNEL")

	PanicChannelID = "1062067299663761428"
	if channel == "" {
	} else {
		PanicChannelID = channel
	}
}
