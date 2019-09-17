package domain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type Fetcher interface {
	GetList(courierID int, minutes int) ([]Coordinates, error)
}

type LocationsFetcher struct {
	client *http.Client
	url    string
}

func NewLocationsFetcher(client *http.Client, url string) *LocationsFetcher {
	return &LocationsFetcher{
		client: client,
		url:    url,
	}
}

func (z LocationsFetcher) GetList(courierID int, minutes int) ([]Coordinates, error) {
	baseURL, err := url.Parse(fmt.Sprintf("http://"+z.url, courierID))
	if err != nil {
		fmt.Println("Malformed URL: ", err.Error())
		return nil, err
	}

	params := url.Values{}
	params.Add("minutes", strconv.Itoa(minutes))
	baseURL.RawQuery = params.Encode()

	res, err := z.client.Get(baseURL.String())

	if err != nil {
		log.Print(err)
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		log.Printf("received http code %d", res.StatusCode)
		return nil, fmt.Errorf("received http code %d", res.StatusCode)
	}

	respBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	if errClose := res.Body.Close(); errClose != nil {
		log.Print(errClose)
	}

	latlongs := []Coordinates{}

	err = json.Unmarshal(respBytes, &latlongs)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	return latlongs, err
}
