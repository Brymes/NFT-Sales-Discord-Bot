package services

import (
	"DIA-NFT-Sales-Bot/config"
	"fmt"
	"math"
	"time"
)

func TrackFloorPrices() {
	ticker := time.NewTicker(120 * time.Second)

	for {
		select {
		case <-ticker.C:
			if config.FloorPriceTrackerAddress == "" {
				continue
			} else {
				response := FloorPriceAPI(config.FloorPriceTrackerAddress, config.FloorPriceTrackerChain)
				rounded := math.Round(response.FloorPrice.FloorPrice*100) / 100
				update := setCurrency(rounded, response.Volume.Collection)
				err := config.DiscordBot.GuildMemberNickname(config.FloorPriceTrackerGuild, "@me", update)
				if err != nil {
					panic(err)
				}
			}
		default:
			continue
		}
	}
}

func setCurrency(rounded float64, collection string) (updated string) {

	switch config.FloorPriceTrackerChain {

	case "Ethereum":
		updated = fmt.Sprintf("%v %s %s", rounded, "ETH", collection)
	case "Astar":
		updated = fmt.Sprintf("%v %s %s", rounded, "ASTR", collection)

	}

	if config.TrackerCurrency == "USD" {

		token := GetAssetQuote("0x0000000000000000000000000000000000000000", config.FloorPriceTrackerChain)
		rounded = math.Round(rounded * token.Price)

		updated = fmt.Sprintf("%s %v %s", "$", rounded, collection)
	}
	return
}
