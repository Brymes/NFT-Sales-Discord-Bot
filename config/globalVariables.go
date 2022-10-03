package config

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
	"sync"
)

var (
	DBClient *gorm.DB

	DiscordBot *discordgo.Session

	// ActiveSales Maps Address to A List of channels
	ActiveSales    = map[string][]string{}
	ActiveSalesMux = &sync.Mutex{}

	// ActiveAllSales Maps Threshold to A List of Channels
	ActiveAllSales     = map[float64][]string{}
	ActiveAllSalesKeys []float64
	ActiveAllSalesMux  = &sync.Mutex{}

	ActiveNftEventWS     = false
	NftEventWSCancelFunc context.CancelFunc

	PanicChannelID string
)

func ShutDownWS() {
	if ActiveNftEventWS {
		NftEventWSCancelFunc()
	}
}
