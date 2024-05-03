package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type VultrRegion struct {
	ID        string   `json:"id"`
	Name      string   `json:"city"`
	Country   string   `json:"country"`
	Continent string   `json:"continent"`
	Options   []string `json:"options"`
}

type VultrResponse struct {
	Regions []VultrRegion `json:"regions"`
}

func getVultrData() VultrResponse {
	url := "https://api.vultr.com/v2/regions"
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var data VultrResponse
	json.Unmarshal(body, &data)

	return data
}

func getVultrStorageRegions(data VultrResponse) map[string]string {
	var regionMap map[string]string = make(map[string]string)
	for _, region := range data.Regions {
		for _, option := range region.Options {
			if option == "block_storage_storage_opt" {
				regionMap[region.ID] = fmt.Sprintf("%s, %s (%s)", region.Name, region.Country, region.Continent)
			}
		}
	}

	return regionMap
}
func getVultrComputeRegions(data VultrResponse) map[string]string {
	var regionMap map[string]string = make(map[string]string)
	for _, region := range data.Regions {
		for _, option := range region.Options {
			if option == "kubernetes" {
				regionMap[region.ID] = fmt.Sprintf("%s, %s (%s)", region.Name, region.Country, region.Continent)
			}
		}
	}

	return regionMap
}

func getVultrRegions() Regions {
	data := getVultrData()

	return Regions{
		Storage: getVultrStorageRegions(data),
		Compute: getVultrComputeRegions(data),
	}
}
