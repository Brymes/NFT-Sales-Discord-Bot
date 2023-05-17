package utils

import (
	log "DIA-NFT-Sales-Bot/debug"
	"fmt"
	"math/big"
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

var BaseBinanceURL = map[string]string{
	"address":     "https://bscscan.com/address",
	"transaction": "https://bscscan.com/tx",
}
var ChainURLMap = map[string]map[string]string{
	"ethereum": BaseEtherScanURL,
	"astar":    BaseAstarURL,
	"solana":   BaseSolanaURL,
	"binance":  BaseBinanceURL,
}

func GetScanLink(linkType, payload, blockchain string) string {
	linkType = strings.ToLower(linkType)
	blockchain = strings.ToLower(blockchain)
	chainUrls := ChainURLMap[blockchain]
	return fmt.Sprintf("<%s/%s>", chainUrls[linkType], payload)
}

var BaseMarketPlaceURL = map[string]string{
	"opensea":                   "https://opensea.io/assets/ethereum",
	"looksrare":                 "https://looksrare.org/collections",
	"x2y2":                      "https://x2y2.io/eth",
	"magiceden":                 "https://explorer.solana.com/address",
	"TofuNFT-Astar":             "https://tofunft.com/nft/astar/",
	"TofuNFT-BinanceSmartChain": "https://tofunft.com/nft/bsc/",
}

func GetMarketPlaceLink(marketPlace, collectionAddress, tokenID string) string {
	marketPlaceURL, match := BaseMarketPlaceURL[strings.ToLower(marketPlace)]
	if match {
		return fmt.Sprintf("<%s/%s/%s>", marketPlaceURL, collectionAddress, tokenID)
	}
	return ""
}

func ConvertDecimalsToCurrency(price int64, decimals int64) *big.Float {
	res := big.NewInt(0)
	fl := big.NewFloat(0)

	// res, isok := res.SetString(decimals, 10)
	// if !isok {
	// 	fmt.Printf("error in exp price %s and decimals %s \n", price, decimals)
	// 	return res
	// }

	res = res.Exp(big.NewInt(10), big.NewInt(decimals), nil)

	resfl, _ := fl.SetString(res.String())

	pricebig := fl.Quo(big.NewFloat(float64(price)), resfl)
	return pricebig
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

func CreateHyperLink(text, url string) string {
	return fmt.Sprintf("[%s](%s)", text, url)
}

func OnErrorPanic(err error, helpText string) {
	if err != nil {
		log.Log.Panicf("%s: \n, %v", helpText, err)
	}
}
