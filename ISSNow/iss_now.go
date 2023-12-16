package ISSNow

import (
	"encoding/json"
	"github.com/kelvins/geocoder"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type IssRaw struct {
	Timestamp   int `json:"timestamp"`
	IssPosition struct {
		Latitude  string `json:"latitude"`
		Longitude string `json:"longitude"`
	} `json:"iss_position"`
	Message string `json:"message"`
}

type Iss struct {
	Latitude  string
	Longitude string
	Location  string
}

var (
	url    = "http://api.open-notify.org/iss-now.json"
	client = http.Client{
		Timeout: time.Second * 10,
	}
)

func GetLocation(apikey string) Iss {
	coords := getISSCoords()

	latitude, _ := strconv.ParseFloat(coords.IssPosition.Latitude, 32)
	longitude, _ := strconv.ParseFloat(coords.IssPosition.Longitude, 32)
	geocoder.ApiKey = apikey
	location := geocoder.Location{
		Latitude:  latitude,
		Longitude: longitude,
	}
	addresses, err := geocoder.GeocodingReverse(location)

	if err != nil {
		log.Fatal("Could not get the addresses: ", err)
	}

	if err != nil {
		log.Fatal(err)
	}
	address := addresses[0]

	iss := Iss{}
	iss.Latitude = coords.IssPosition.Latitude
	iss.Longitude = coords.IssPosition.Longitude
	iss.Location = address.FormatAddress()

	return iss
}

func getISSCoords() IssRaw {
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

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	iss := IssRaw{}
	jsonErr := json.Unmarshal(body, &iss)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return iss
}
