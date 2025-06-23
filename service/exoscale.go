package service

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getExoscaleStorageRegions() map[string]string {
	url := "https://www.exoscale.com/datacenters/"
	doc, err := get(url)
	if err != nil || doc == nil {
		// fmt.Printf("[Exoscale] Error fetching or parsing regions: %v\n", err)
		return map[string]string{}
	}

	var regionMap map[string]string = make(map[string]string)

	// Find the datacenters div and parse article elements
	doc.Find("div.datacenters article").Each(func(i int, article *goquery.Selection) {
		// Extract locality and region code from spans within h2
		locality := strings.TrimSpace(article.Find("span.datacenters-locality").Text())
		regionCode := strings.TrimSpace(article.Find("span.datacenters-name").Text())

		if locality != "" && regionCode != "" {
			regionMap[regionCode] = fmt.Sprintf("%s - %s", locality, regionCode)
		}
	})

	return regionMap
}

func GetExoscaleRegions() Regions {
	regions := getExoscaleStorageRegions()
	return Regions{
		Storage: regions,
		Compute: regions,
	}
}
