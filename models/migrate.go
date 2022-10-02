package models

import (
	"DIA-NFT-Sales-Bot/config"
	"DIA-NFT-Sales-Bot/utils"
	"log"
	"sort"
)

func InitMigrations() {
	err := config.DBClient.AutoMigrate(&Subscriptions{})
	if err != nil {
		log.Println("Error performing Database Migrations")
		log.Fatalln(err)
	}
}

func LoadCurrentSubscriptions() bool {
	subscriptions := Subscriptions{}.LoadAllSubscriptions()

	for _, subscription := range subscriptions {
		switch subscription.Command {
		case "sales":
			go func() {
				config.ActiveSalesMux.Lock()

				subscribedChannels := config.ActiveSales[subscription.Address.String]
				subscribedChannels = append(subscribedChannels, subscription.ChannelID.String)
				config.ActiveSales[subscription.Address.String] = utils.RemoveArrayDuplicates(subscribedChannels)

				config.ActiveSalesMux.Unlock()
			}()
		case "all_sales":
			go func() {
				config.ActiveAllSalesMux.Lock()

				subscribedChannels := config.ActiveAllSales[subscription.Threshold]
				subscribedChannels = append(subscribedChannels, subscription.ChannelID.String)

				config.ActiveAllSales[subscription.Threshold] = utils.RemoveArrayDuplicates(subscribedChannels)

				config.ActiveAllSalesKeys = make([]float64, 0, len(subscribedChannels))

				for k := range config.ActiveAllSales {
					config.ActiveAllSalesKeys = append(config.ActiveAllSalesKeys, k)
				}

				sort.Float64s(config.ActiveAllSalesKeys)

				config.ActiveAllSalesMux.Unlock()
			}()
		}
	}
	if len(subscriptions) > 0 {
		return true
	}
	return false
}
