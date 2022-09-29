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
	ActiveSalesMux = &sync.RWMutex{}

	// ActiveAllSales Maps Threshold to A List of Channels
	ActiveAllSales     = map[float64][]string{}
	ActiveAllSalesKeys []float64
	ActiveAllSalesMux  = &sync.RWMutex{}

	ActiveNftEventWS     = false
	NftEventWSCancelFunc context.CancelFunc

	PanicChannelID string
)
