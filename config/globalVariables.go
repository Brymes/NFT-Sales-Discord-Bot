package config

import (
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
	"sync"
)

var (
	DBClient *gorm.DB

	DiscordBot *discordgo.Session

	ActiveSales    = map[string][]string{}
	ActiveSalesMux = &sync.RWMutex{}

	ActiveAllSales     = map[float64][]string{}
	ActiveAllSalesKeys []float64
	ActiveAllSalesReversed []float64
	ActiveAllSalesMux = &sync.RWMutex{}


	ActiveNftEventWS = false

	PanicChannelID string
)