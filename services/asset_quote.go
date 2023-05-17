package services

import (
	"DIA-NFT-Sales-Bot/utils"
	"encoding/json"
	"fmt"
	"time"
)

type AssetPricePayload struct {
	Symbol             string    `json:"Symbol"`
	Name               string    `json:"Name"`
	Address            string    `json:"Address"`
	Blockchain         string    `json:"Blockchain"`
	Price              float64   `json:"Price"`
	PriceYesterday     float64   `json:"PriceYesterday"`
	VolumeYesterdayUSD float64   `json:"VolumeYesterdayUSD"`
	Time               time.Time `json:"Time"`
	Source             string    `json:"Source"`
	Signature          string    `json:"Signature"`
}

func GetAssetQuote(contractAddress, blockchain string) (response AssetPricePayload) {
	response = AssetPricePayload{}
	url := fmt.Sprintf("https://api.diadata.org/v1/assetQuotation/%s/%s", blockchain, contractAddress)

	bodyCloser, body := MakeRequest(url)
	err := json.Unmarshal(body, &response)

	bodyCloser.Close()

	utils.OnErrorPanic(err, "Error with Getting Asset Quote")

	return response
}
