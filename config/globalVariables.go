package config

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
	"os"
	"sync"
)

type SubscriptionChannelArray string

var (
	DBClient *gorm.DB

	DiscordBot *discordgo.Session

	// ActiveSales Maps Address to A List of channels
	ActiveSales    = map[string]map[string][]SubscriptionChannelArray{}
	ActiveSalesMux = &sync.Mutex{}

	// ActiveAllSales Maps Threshold to A List of Channels
	ActiveAllSales = map[float64]map[string][]SubscriptionChannelArray{}
	ActiveAllSalesKeys []float64
	ActiveAllSalesMux = &sync.Mutex{}

	ActiveNftEventWS     = false
	NftEventWSCancelFunc context.CancelFunc

	// PanicChannelID Variable to Hold ChannelID to forward all errors
	PanicChannelID string
)

func ShutDownWS() {
	if ActiveNftEventWS {
		NftEventWSCancelFunc()
	}
}

func InitPanicChannel() {
	channel := os.Getenv("PANIC_CHANNEL")

	if channel == "" {
		PanicChannelID = "1025726821733515314"
	} else {
		PanicChannelID = channel
	}
}
