package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
		jsonContent := `{
			"storage": {
				"ap-south-1": "Singapore - ap-south-1",
				"br-gru-1": "São Paulo (Brazil) - br-gru-1",
				"es-mad-1": "Madrid (Spain) - es-mad-1",
				"eu-central-1": "Frankfurt (Germany) - eu-central-1",
				"fr-par-1": "Paris (France) - fr-par-1",
				"id-cgk-1": "Jakarta (Indonesia) - id-cgk-1",
				"in-maa-1": "Chennai (India) - in-maa-1",
				"it-mil-1": "Milan (Italy) - it-mil-1",
				"jp-osa-1": "Osaka (Japan) - jp-osa-1",
				"nl-ams-1": "Amsterdam (Netherlands) - nl-ams-1",
				"se-sto-1": "Stockholm (Sweden) - se-sto-1",
				"us-east-1": "Newark, NJ (USA) - us-east-1",
				"us-iad-1": "Washington, DC (USA) - us-iad-1",
				"us-lax-1": "Los Angeles, CA (USA) - us-lax-1",
				"us-mia-1": "Miami, FL (USA) - us-mia-1",
				"us-ord-1": "Chicago, IL (USA) - us-ord-1",
				"us-sea-1": "Seattle, WA (USA) - us-sea-1",
				"us-southeast-1": "Atlanta, GA (USA) - us-southeast-1"
			},
			"compute": {
				"ap-northeast": "Tokyo, JP - ap-northeast",
				"ap-south": "Singapore, SG - ap-south",
				"ap-southeast": "Sydney, AU - ap-southeast",
				"ap-west": "Mumbai, IN - ap-west",
				"br-gru": "Sao Paulo, BR - br-gru",
				"ca-central": "Toronto, CA - ca-central",
				"es-mad": "Madrid, ES - es-mad",
				"eu-central": "Frankfurt, DE - eu-central",
				"eu-west": "London, UK - eu-west",
				"fr-par": "Paris, FR - fr-par",
				"id-cgk": "Jakarta, ID - id-cgk",
				"in-maa": "Chennai, IN - in-maa",
				"it-mil": "Milan, IT - it-mil",
				"jp-osa": "Osaka, JP - jp-osa",
				"nl-ams": "Amsterdam, NL - nl-ams",
				"se-sto": "Stockholm, SE - se-sto",
				"us-central": "Dallas, TX - us-central",
				"us-east": "Newark, NJ - us-east",
				"us-iad": "Washington, DC - us-iad",
				"us-lax": "Los Angeles, CA - us-lax",
				"us-mia": "Miami, FL - us-mia",
				"us-ord": "Chicago, IL - us-ord",
				"us-sea": "Seattle, WA - us-sea",
				"us-southeast": "Atlanta, GA - us-southeast",
				"us-west": "Fremont, CA - us-west"
			}
		}`

		// json has "storage" and "compute" keys
		var regions map[string]map[string]string
		json.Unmarshal([]byte(jsonContent), &regions)

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
