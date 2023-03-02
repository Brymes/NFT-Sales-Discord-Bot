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
				update := fmt.Sprintf("%v %s", rounded, response.Volume.Collection)
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
