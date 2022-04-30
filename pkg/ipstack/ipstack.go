package ipstack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var (
	uri = "http://api.ipstack.com/"
	key = "9796fdad009d246b4d19505bcea99944"
)

type info struct {
	IPType      string `json:"type"`
	Continent   string `json:"continent_name"`
	CountryCode string `json:"country_code"`
	CountryName string `json:"country_name"`
	Region      string `json:"region_name"`
	City        string `json:"city"`
}

func GetInfo(ip string) (string, error) {
	var output string

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req := &http.Request{
		Method: http.MethodGet,
	}

	req.URL, _ = url.Parse(uri + ip + "?access_key=" + key)

	resp, err := client.Do(req)
	if err != nil {
		return output, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return output, err
	}

	info := &info{}

	if err := json.Unmarshal(data, info); err != nil {
		return output, err
	}

	output += fmt.Sprintf(
		"IP: %s\nType: %s\nContinent: %s\nCountry code: %s\nCountry name: %s\nRegion: %s\nCity: %s",
		ip, info.IPType, info.Continent, info.CountryCode, info.CountryName, info.Region, info.City,
	)
	return output, nil
}
