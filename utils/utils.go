package utils

import (
	"fmt"
	"math"
	"strings"
)

var BaseEtherScanURL = map[string]string{
	"address":     "https://etherscan.io/address",
	"transaction": "https://etherscan.io/tx",
}
var BaseSolanaURL = map[string]string{
	"address":     "https://explorer.solana.com/address",
	"transaction": "https://explorer.solana.com/tx",
}
var BaseAstarURL = map[string]string{
	"address":     "https://astar.subscan.io/account",
	"transaction": "https://astar.subscan.io/extrinsic",
}
var ChainURLMap = map[string]map[string]string{
	"ethereum": BaseEtherScanURL,
	"astar":    BaseAstarURL,
	"solana":   BaseSolanaURL,
}

func GetScanLink(linkType, payload, blockchain string) string {
	linkType = strings.ToLower(linkType)
	blockchain = strings.ToLower(blockchain)
	chainUrls := ChainURLMap[blockchain]
	return fmt.Sprintf("%s/%s", chainUrls[linkType], payload)
}

var BaseMarketPlaceURL = map[string]string{
	"opensea":   "https://opensea.io/assets/ethereum",
	"looksrare": "https://looksrare.org/collections",
	"x2y2":      "https://x2y2.io/eth",
	"magiceden": "https://explorer.solana.com/address",
}

func GetMarketPlaceLink(marketPlace, collectionAddress, tokenID string) string {
	marketPlaceURL, match := BaseMarketPlaceURL[strings.ToLower(marketPlace)]
	if match {
		return fmt.Sprintf("%s/%s/%s", marketPlaceURL, collectionAddress, tokenID)
	}
	return ""
}

func ConvertDecimalsToCurrency(price int64, decimals int) float64 {
	res := math.Pow(10, float64(decimals))
	priceInEth := float64(price) / res
	return priceInEth
}

func RemoveArrayDuplicates(arr []string) []string {
	occurred := map[string]bool{}
	var result []string

	for e := range arr {

		// check if already the mapped
		// variable is set to true or not
		if occurred[arr[e]] != true {
			occurred[arr[e]] = true

			// Append to result slice.
			result = append(result, arr[e])
		}
	}

	return result
}
