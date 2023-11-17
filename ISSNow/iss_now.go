package ISSNow

import (
	"io"
	"log"
	"net/http"
	"time"
)

type iss_coords struct {
	latitude  float64 `json:"latitude"`
	longitude float64 `json:"longitude"`
}

type iss struct {
	coords   iss_coords
	location string
}

var (
	url = "http://api.open-notify.org/iss-now.json"
)

func GetISSCoords() iss_coords {
	client := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := io.ReadAll(res.Body)

	if readErr != nil {
		log.Fatal(readErr)
	}

	coords := iss_coords{}

}
