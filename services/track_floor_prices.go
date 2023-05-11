package services

import (
	"DIA-NFT-Sales-Bot/config"
	"fmt"
	"math"
	"time"
)

func TrackFloorPrices() {
	ticker := time.NewTicker(20 * time.Second)

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
	//Get Price for WETH to USD
	if config.TrackerCurrency == "USD" {
		wethPrice := GetAssetQuote("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2", "Ethereum")

		rounded = math.Round(rounded*wethPrice.Price*10) / 10

		updated = fmt.Sprintf("%s %v %s", "$", rounded, collection)
	} else {
		updated = fmt.Sprintf("%v %s %s", rounded, "ETH", collection)
	}
	return
}
