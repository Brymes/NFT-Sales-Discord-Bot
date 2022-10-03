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
