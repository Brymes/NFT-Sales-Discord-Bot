package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"
)

func LastTradesAPI(address, blockchain string) (response []LastTradesAPIResponse) {
	var (
		url        string
		bodyCloser io.ReadCloser
		body       []byte
		err        error
	)
	response = []LastTradesAPIResponse{}

	url = fmt.Sprintf("https://api.diadata.org/v1/NFTTradesCollection/%s/%s", blockchain, address)

	bodyCloser, body = MakeRequest(url)
	err = json.Unmarshal(body, &response)
	if err != nil {
		errorRes := errors.New("Error reading response from LastTrades API" + err.Error() + "\n\n" + address + blockchain)
		panic(errorRes)
	}

	bodyCloser.Close()

	return
}

type LastTradesAPIResponse struct {
	Name        string    `json:"Name"`
	Price       float64   `json:"Price"`
	NFTid       string    `json:"NFTid"`
	FromAddress string    `json:"FromAddress"`
	ToAddress   string    `json:"ToAddress"`
	BundleSale  bool      `json:"BundleSale"`
	BlockNumber int       `json:"BlockNumber"`
	Timestamp   time.Time `json:"Timestamp"`
	TxHash      string    `json:"TxHash"`
	Exchange    string    `json:"Exchange"`
	Currency    struct {
		Symbol     string `json:"Symbol"`
		Name       string `json:"Name"`
		Address    string `json:"Address"`
		Decimals   int    `json:"Decimals"`
		Blockchain string `json:"Blockchain"`
	} `json:"Currency"`
}
