package services

import (
	"DIA-NFT-Sales-Bot/config"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"time"
)

type NFTEvent struct {
	Error    string `json:"Error"`
	Response struct {
		NFT struct {
			NFTClass struct {
				Address      string `json:"Address"`
				Symbol       string `json:"Symbol"`
				Name         string `json:"Name"`
				Blockchain   string `json:"Blockchain"`
				ContractType string `json:"ContractType"`
				Category     string `json:"Category"`
			} `json:"NFTClass"`
			TokenID string `json:"TokenID"`
		} `json:"NFT"`
		Price       int64  `json:"Price"`
		FromAddress string `json:"FromAddress"`
		ToAddress   string `json:"ToAddress"`
		Currency    struct {
			Symbol     string `json:"Symbol"`
			Name       string `json:"Name"`
			Address    string `json:"Address"`
			Decimals   int    `json:"Decimals"`
			Blockchain string `json:"Blockchain"`
		} `json:"Currency"`
		BundleSale  bool      `json:"BundleSale"`
		BlockNumber int       `json:"BlockNumber"`
		Timestamp   time.Time `json:"Timestamp"`
		TxHash      string    `json:"TxHash"`
		Exchange    string    `json:"Exchange"`
	} `json:"Response"`
}

func ConnectToService(logger *log.Logger) {

	u := url.URL{Scheme: "wss", Host: "api.diadata.org", Path: "/ws/nft"}

	logger.Printf("connecting to %s", u.String())
	c, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		config.ActiveNftEventWS = false
		logger.Printf("handshake failed with status %d", resp.StatusCode)
		logger.Fatal("dial:", err)
	}

	//When the program closes close the connection
	defer c.Close()

	err = c.WriteJSON(map[string]string{"Channel": "nftsales"})
	if err != nil {
		logger.Println("ks1")
		logger.Println(err)
		config.ActiveNftEventWS = false
	}

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			event, jsonPayload := NFTEvent{}, map[string]interface{}{}

			err = c.ReadJSON(&jsonPayload)
			if err != nil {
				logger.Println("Error Reading response", err)
			}

			check := jsonPayload["Response"]

			switch check {

			case "alive":
				config.ActiveNftEventWS = true
				break
			case "subscibed to nftsales":
				config.ActiveNftEventWS = true
			default:
				// convert map to json
				jsonString, _ := json.Marshal(jsonPayload)
				fmt.Println(string(jsonString))
				// convert json to struct
				err := json.Unmarshal(jsonString, &event)
				if err != nil {
					logger.Println("Error Unmarshalling to Struct")
					logger.Println(err)
				}
				SalesController(event)
			}
		}
	}()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			config.ActiveNftEventWS = false
			return
		case t := <-ticker.C:
			err = c.WriteJSON(map[string]string{"Channel": "ping"})

			if err != nil {
				logger.Println(t)
				logger.Println("write:", err)
				config.ActiveNftEventWS = false
				return
			}
		}
	}
}
