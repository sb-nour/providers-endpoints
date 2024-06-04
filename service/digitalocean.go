package service

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// getDigitalOceanSpacesRegions retrieves the regions for DigitalOcean Spaces.
// It makes a GET request to the DigitalOcean Spaces availability URL and parses the HTML response to extract the regions.
// The regions are then translated using the translateRegions function.
// Returns a map of region names and their corresponding values.
func getDigitalOceanSpacesRegions() map[string]string {
	url := "https://docs.digitalocean.com/products/spaces/details/availability/"
	doc, _ := get(url)

	var regions []string
	doc.Find("thead th").Each(func(index int, th *goquery.Selection) {
		// Check if the corresponding <td> has the "◆" symbol
		if td := doc.Find("tbody tr").Children().Eq(index); td.Text() == "◆" {
			regions = append(regions, th.Text())
		}
	})

	return translateRegions(regions)
}

func getDigitalOceanDropletRegions() map[string]string {
	url := "https://docs.digitalocean.com/platform/regional-availability/"
	doc, _ := get(url)

	var regionMap = make(map[string]string)
	doc.Find("table").Each(func(index int, table *goquery.Selection) {
		if table.Find("thead th").First().Text() == "Datacenter" {
			table.Find("tbody tr").Each(func(index int, tr *goquery.Selection) {
				regionCode := strings.ToLower(tr.Children().Eq(2).Text())
				regionName := fmt.Sprintf("%s - %s", tr.Children().Eq(1).Text(), regionCode)
				regionMap[regionCode] = regionName
			})
		}
	})
	return regionMap
}

func GetDigitalOceanRegions() Regions {
	return Regions{
		Storage: getDigitalOceanSpacesRegions(),
		Compute: getDigitalOceanDropletRegions(),
	}
}
