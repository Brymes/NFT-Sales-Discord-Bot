package config

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
	"os"
	"sync"
)

var (
	DBClient *gorm.DB

	DiscordBot *discordgo.Session

	// ActiveSales Maps Address to A List of channels
	ActiveSales    = map[string]map[string][]string{}
	ActiveSalesMux = &sync.Mutex{}

	// ActiveAllSales Maps Threshold to A List of Channels
	ActiveAllSales     = map[float64]map[string][]string{}
	ActiveAllSalesKeys []float64
	ActiveAllSalesMux  = &sync.Mutex{}

	ActiveNftEventWS     = false
	NftEventWSCancelFunc context.CancelFunc

	ActiveSalesInfoBot = map[string]map[string]string{}
	ActiveSalesInfoMux = &sync.Mutex{}

	// PanicChannelID Variable to Hold ChannelID to forward all errors
	PanicChannelID string

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

	PanicChannelID = "1025726821733515314"
	if channel == "" {
	} else {
		PanicChannelID = channel
	}
}
