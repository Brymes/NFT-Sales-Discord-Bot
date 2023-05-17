package models

import (
	"DIA-NFT-Sales-Bot/config"
	log "DIA-NFT-Sales-Bot/debug"
	"DIA-NFT-Sales-Bot/utils"
	"math/big"
	"sort"
	"strings"
)

func InitMigrations() {
	err := config.DBClient.AutoMigrate(&Subscriptions{}, &ConfigModel{})
	if err != nil {
		log.Log.Println("Error performing Database Migrations")
		log.Log.Fatalln(err)
	}
}

func LoadCurrentSubscriptions() bool {
	subscriptions := Subscriptions{}.LoadAllSubscriptions()
	cm := ConfigModel{}
	cm.GetConfig()
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
			threshold := big.NewFloat(0)
			threshold, _ = threshold.SetString(subscription.Threshold.String)

			subscribedChannels := config.ActiveAllSales[threshold][subscription.Blockchain]
			subscribedChannels = append(subscribedChannels, subscription.ChannelID.String)

			data := config.ActiveAllSales[threshold]
			if len(data) == 0 {
				data = map[string][]string{}
			}
			data[subscription.Blockchain] = utils.RemoveArrayDuplicates(subscribedChannels)
			config.ActiveAllSales[threshold] = data

			config.ActiveAllSalesKeys = make([]*big.Float, 0, len(subscribedChannels))

			for k := range config.ActiveAllSales {
				config.ActiveAllSalesKeys = append(config.ActiveAllSalesKeys, k)
			}

			// sort.Float64s(config.ActiveAllSalesKeys)

			sort.Slice(config.ActiveAllSalesKeys, func(a, b int) bool {
				return config.ActiveAllSalesKeys[a].Cmp(config.ActiveAllSalesKeys[b]) > 0
			})

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
