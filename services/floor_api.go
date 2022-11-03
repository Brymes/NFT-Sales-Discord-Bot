package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type FloorPriceResponse struct {
	FloorPrice float64 `json:"Floor_Price"`
}
type MovingAverageAPIResponse struct {
	MovingAverageFloorPrice float64 `json:"Moving_Average_Floor_Price"`
}
type VolumeAPIResponse struct {
	Collection   string    `json:"Collection"`
	Floor        float64   `json:"Floor"`
	Volume       float64   `json:"Volume"`
	Trades       int       `json:"Trades"`
	FloorChange  float64   `json:"FloorChange"`
	VolumeChange float64   `json:"VolumeChange"`
	TradesChange float64   `json:"TradesChange"`
	Address      string    `json:"Address"`
	Blockchain   string    `json:"Blockchain"`
	Time         time.Time `json:"Time"`
	Source       string    `json:"Source"`
	Exchanges    []struct {
		Exchange  string  `json:"Exchange"`
		NumTrades int     `json:"NumTrades"`
		Volume    float64 `json:"Volume"`
	} `json:"Exchanges"`
}

type Floor struct {
	Volume     VolumeAPIResponse
	MA         MovingAverageAPIResponse
	FloorPrice FloorPriceResponse
}

func FloorPriceAPI(contractAddress, blockchain string) (response Floor) {
	var (
		urls       = ParseURLs(contractAddress, blockchain)
		bodyCloser io.ReadCloser
		body       []byte
		err        error
	)

	//x := map[string]interface{}{}

	response = Floor{}

	for key, url := range urls {
		switch key {

		case "volume":
			bodyCloser, body = MakeRequest(url)
			err = json.Unmarshal(body, &response.Volume)
		case "movingAverage":
			bodyCloser, body = MakeRequest(url)
			err = json.Unmarshal(body, &response.MA)
		case "floorPrice":
			bodyCloser, body = MakeRequest(url)
			err = json.Unmarshal(body, &response.FloorPrice)
		}
		bodyCloser.Close()
		if err != nil {
			errorRes := errors.New("Error reading response from Floor Price API" + err.Error())
			panic(errorRes)
		}
	}

	return response
}

func MakeRequest(url string) (io.ReadCloser, []byte) {
	var (
		errorRes error
		method   = "GET"
		client   = &http.Client{}
	)

	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Println(err)
		errorRes = errors.New("Error whilst setting up Communication with Floor Price API" + err.Error())
		panic(errorRes)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		errorRes = errors.New("Error Communicating with Floor Price API" + err.Error())
		panic(errorRes)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		errorRes = errors.New("Error reading response from Floor Price API" + err.Error())
		panic(errorRes)
	}
	return res.Body, body

}

func ParseURLs(address, blockchain string) map[string]string {
	urls := map[string]string{}

	urls["floorPrice"] = fmt.Sprintf("https://api.diadata.org/v1/NFTFloor/%s/%s", blockchain, address)
	urls["movingAverage"] = fmt.Sprintf("https://api.diadata.org/v1/NFTFloorMA/%s/%s?floorWindow=86400", blockchain, address)
	urls["volume"] = fmt.Sprintf("https://api.diadata.org/v1/NFTVolume/%s/%s?starttime=%v&endtime=%v", blockchain, address, time.Now().Add(-(24 * time.Hour)).Unix(), time.Now().Unix())

	return urls
}
