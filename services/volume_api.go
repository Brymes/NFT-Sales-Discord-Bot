package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"
)

func VolumeAPI(address, blockchain string) (response VolumeAPIResponse) {
	var (
		url        string
		bodyCloser io.ReadCloser
		body       []byte
		err        error
	)
	response = VolumeAPIResponse{}

	url = fmt.Sprintf("https://api.diadata.org/v1/NFTVolume/%s/%s?starttime=%v&endtime=%v", blockchain, address, time.Now().Add(-(24 * time.Hour)).Unix(), time.Now().Unix())

	bodyCloser, body = MakeRequest(url)
	err = json.Unmarshal(body, &response)
	if err != nil {
		errorRes := errors.New("Error reading response from Volume API" + err.Error() + "\n\n" + address + blockchain)
		panic(errorRes)
	}

	bodyCloser.Close()

	return
}
