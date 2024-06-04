package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type LinodeRegion struct {
	ID      string   `json:"id"`
	Label   string   `json:"label"`
	Country string   `json:"country"`
	Options []string `json:"capabilities"`
}

type LinodeResponse struct {
	Regions []LinodeRegion `json:"data"`
}

func getLinodeStorageRegions() (map[string]string, error) {
	url := "https://www.linode.com/docs/products/storage/object-storage/"
	doc, err := get(url)
	if err != nil {
		return nil, err
	}

	var regionMap map[string]string = make(map[string]string)
	doc.Find("table").Each(func(i int, table *goquery.Selection) {
		// if table doesn't have more than 2 rows, return
		if table.Find("tbody tr").Length() < 5 {
			return
		}

		if table.Find("thead th").Length() == 2 && table.Find("thead th").Eq(0).Text() == "Data Center" {

			table.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
				regionCode := row.Find("td").Eq(1).Text()
				regionName := strings.ReplaceAll(row.Find("td").Eq(0).Text(), "*", "")
				regionMap[regionCode] = fmt.Sprintf("%s - %s", regionName, regionCode)
			})
		}
	})

	return regionMap, nil
}

func getLinodeData() LinodeResponse {
	url := "https://api.linode.com/v4/regions"
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var data LinodeResponse
	json.Unmarshal(body, &data)

	return data
}

func getLinodeComputeRegions(data LinodeResponse) map[string]string {
	var regionMap map[string]string = make(map[string]string)
	for _, region := range data.Regions {
		for _, option := range region.Options {
			if option == "Linodes" {
				regionMap[region.ID] = fmt.Sprintf("%s - %s", region.Label, region.ID)
			}
		}
	}

	return regionMap
}

func GetLinodeRegions() Regions {
	storageRegions, err := getLinodeStorageRegions()
	if err != nil {
		// fmt.Printf("Error: %v", err)
		// Load the regions from the local file ./linode_fallback.json
		filePath := "./service/linode_fallback.json"
		jsonContent, err := os.ReadFile(filePath)
		if err != nil {
			if debugging {
				fmt.Printf("Error reading file: %v", err)
			}
			return Regions{}
		}

		// json has "storage" and "compute" keys
		var regions map[string]map[string]string
		json.Unmarshal(jsonContent, &regions)

		return Regions{
			Storage: regions["storage"],
			Compute: regions["compute"],
		}
	}
	return Regions{
		Storage: storageRegions,
		Compute: getLinodeComputeRegions(getLinodeData()),
	}
}
