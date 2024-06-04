package service

import (
	"encoding/json"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getUpcloudStorageRegions(doc *goquery.Document) map[string]string {
	var regionMap map[string]string = make(map[string]string)

	doc.Find(".accordion").First().Find(".accordion-item").Each(func(i int, item *goquery.Selection) {
		// Check if item has an <li> tag with "Object Storage" text
		if strings.Contains(item.Find("li").Text(), "Object Storage") {

			regionCode := strings.ToLower(item.Find("button h3").First().Text())
			regionSplit := strings.Split(item.Find("button .location").First().Text(), ", ")
			regionName := regionSplit[1] + " - " + regionSplit[0] + " - " + regionCode

			regionMap[regionCode] = regionName
		}
	})

	return regionMap
}
func getUpcloudComputeRegions(doc *goquery.Document) map[string]string {

	var regionMap map[string]string = make(map[string]string)

	doc.Find(".accordion").First().Find(".accordion-item").Each(func(i int, item *goquery.Selection) {
		// Check if item has an <li> tag with "Object Storage" text
		if strings.Contains(item.Find("li").Text(), "Cloud Servers") {

			regionCode := strings.ToLower(item.Find("button h3").First().Text())
			regionSplit := strings.Split(item.Find("button .location").First().Text(), ", ")
			regionName := regionSplit[1] + " - " + regionSplit[0] + " - " + regionCode

			regionMap[regionCode] = regionName
		}
	})

	return regionMap
}

func GetUpcloudRegions() Regions {
	doc, err := get("https://upcloud.com/data-centres")
	if err != nil {
		// fmt.Printf("Error: %v", err)
		// Load the regions from the local file ./upcloud_fallback.json
		jsonContent := []byte(`{
  "storage": {
    "au-syd1": "Sydney, Australia - AU SYD1",
    "de-fra1": "Frankfurt, Germany - DE FRA1",
    "fi-hel2": "Helsinki, Finland - FI HEL2",
    "es-mad1": "Madrid, Spain - ES MAD1",
    "nl-ams1": "Amsterdam, Netherlands - NL AMS1",
    "pl-waw1": "Warsaw, Poland - PL WAW1",
    "sg-sin1": "Singapore - SG SIN1",
    "uk-lon1": "London, UK - UK LON1",
    "us-chi1": "Chicago, USA - US CHI1",
    "us-nyc1": "New York, USA - US NYC1",
    "us-sjo1": "San Jose, USA - US SJO1"
  },
  "compute": {
    "au-syd1": "Sydney, Australia - AU SYD1",
    "de-fra1": "Frankfurt, Germany - DE FRA1",
    "fi-hel1": "Helsinki, Finland - FI HEL1",
    "fi-hel2": "Helsinki, Finland - FI HEL2",
    "es-mad1": "Madrid, Spain - ES MAD1",
    "nl-ams1": "Amsterdam, Netherlands - NL AMS1",
    "pl-waw1": "Warsaw, Poland - PL WAW1",
    "se-sto1": "Stockholm, Sweden - SE STO1",
    "sg-sin1": "Singapore - SG SIN1",
    "uk-lon1": "London, UK - UK LON1",
    "us-chi1": "Chicago, USA - US CHI1",
    "us-nyc1": "New York, USA - US NYC1",
    "us-sjo1": "San Jose, USA - US SJO1"
  }
}
`)

		// json has "storage" and "compute" keys
		var regions map[string]map[string]string
		json.Unmarshal(jsonContent, &regions)

		return Regions{
			Storage: regions["storage"],
			Compute: regions["compute"],
		}
	}
	storageRegions := getUpcloudStorageRegions(doc)
	computeRegions := getUpcloudComputeRegions(doc)

	return Regions{
		Storage: storageRegions,
		Compute: computeRegions,
	}
}
