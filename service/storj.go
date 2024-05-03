package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func getStorjStorageRegions() map[string]string {
	url := "https://us1.storj.io/api/v0/config"

	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	regionMap := make(map[string]string)
	var data map[string][]map[string]string
	json.Unmarshal(body, &data)

	for _, satellite := range data["partneredSatellites"] {
		regionMap[satellite["name"]] = satellite["name"]
	}

	return regionMap
}

func GetStorjRegions() Regions {
	return Regions{
		Storage: getStorjStorageRegions(),
	}
}
