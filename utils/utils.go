package utils

import (
	"fmt"
	"log"
	"math"
	"runtime/debug"
	"strings"
)

var BaseEtherScanURL = map[string]string{
	"address":     "https://etherscan.io/address",
	"transaction": "https://etherscan.io/tx",
}

func GetEtherScanLink(linkType, payload string) string {
	linkType = strings.ToLower(linkType)
	return fmt.Sprintf("%s/%s", BaseEtherScanURL[linkType], payload)
}

var BaseMarketPlaceURL = map[string]string{
	"opensea":   "https://opensea.io/assets/ethereum",
	"looksrare": "https://looksrare.org/collections",
	"x2y2":      "https://x2y2.io/eth",
}

func GetMarketPlaceLink(marketPlace, collectionAddress, tokenID string) string {
	marketPlace = strings.ToLower(marketPlace)
	return fmt.Sprintf("%s/%s/%s", BaseMarketPlaceURL[marketPlace], collectionAddress, tokenID)
}
func ConvertDecimalsToEth(price int64, decimals int) float64 {
	res := math.Pow(10, float64(decimals))
	priceInEth := float64(price) / res
	return priceInEth
}

func HandlePanic(logger *log.Logger, customMessage string) {
	//Notify Admin of any uncaught errors
	if err := recover(); err != nil {
		logger.Println(err)

		if customMessage != "" {
			logger.Println(customMessage)
		} else {
			logger.Println(string(debug.Stack()))
		}
	}
}
