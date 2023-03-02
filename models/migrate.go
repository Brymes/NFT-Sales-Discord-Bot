package models

import (
	"DIA-NFT-Sales-Bot/config"
	"DIA-NFT-Sales-Bot/utils"
	"log"
	"sort"
	"strings"
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
	res := false

	for _, subscription := range subscriptions {
		switch subscription.Command {
		case "sales":
			res = true
			config.ActiveSalesMux.Lock()

			subscribedChannels := config.ActiveSales[strings.ToUpper(subscription.Address.String)][subscription.Blockchain]
			subscribedChannels = append(subscribedChannels, subscription.ChannelID.String)
			data := config.ActiveSales[strings.ToUpper(subscription.Address.String)]
			if len(data) == 0 {
				data = make(map[string][]string)
			}

			data[subscription.Blockchain] = utils.RemoveArrayDuplicates(subscribedChannels)
			config.ActiveSales[strings.ToUpper(subscription.Address.String)] = data
			config.ActiveSalesMux.Unlock()

		case "all_sales":
			res = true

			config.ActiveAllSalesMux.Lock()

			subscribedChannels := config.ActiveAllSales[subscription.Threshold][subscription.Blockchain]
			subscribedChannels = append(subscribedChannels, subscription.ChannelID.String)

			data := config.ActiveAllSales[subscription.Threshold]
			if len(data) == 0 {
				data = map[string][]string{}
			}
			data[subscription.Blockchain] = utils.RemoveArrayDuplicates(subscribedChannels)
			config.ActiveAllSales[subscription.Threshold] = data

			config.ActiveAllSalesKeys = make([]float64, 0, len(subscribedChannels))

			for k := range config.ActiveAllSales {
				config.ActiveAllSalesKeys = append(config.ActiveAllSalesKeys, k)
			}

			sort.Float64s(config.ActiveAllSalesKeys)

			config.ActiveAllSalesMux.Unlock()

		case "set_up_info":
			res = true

			config.ActiveSalesInfoMux.Lock()

			config.ActiveSalesInfoBot[subscription.ChannelID.String] = map[string]string{
				"address": subscription.Address.String, "blockchain": subscription.Blockchain,
			}

			config.ActiveSalesInfoMux.Unlock()
		case "track_floor_price":
			config.FloorPriceTrackerAddress, config.FloorPriceTrackerChain, config.FloorPriceTrackerGuild = subscription.Address.String, subscription.Blockchain, subscription.ChannelID.String

		}
	}

	return res
}
