package models

import (
	"DIA-NFT-Sales-Bot/config"
	"DIA-NFT-Sales-Bot/utils"
	"log"
	"math/big"
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

	for _, subscription := range subscriptions {
		switch subscription.Command {
		case "sales":
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

			config.ActiveAllSalesMux.Lock()
			threshold := big.NewFloat(0)
			threshold, _ = threshold.SetString(subscription.Threshold)

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

			config.ActiveSalesInfoMux.Lock()

			config.ActiveSalesInfoBot[subscription.ChannelID.String] = map[string]string{
				"address": subscription.Address.String, "blockchain": subscription.Blockchain,
			}

			config.ActiveSalesInfoMux.Unlock()
		}
	}
	if len(subscriptions) > 0 {
		return true
	}
	return false
}
