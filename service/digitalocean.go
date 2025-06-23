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

	// Get all <th> from <thead> (region codes), skipping the first ("Product")
	headers := []string{}
	doc.Find("thead th").Each(func(i int, th *goquery.Selection) {
		if i > 0 { // skip "Product"
			headers = append(headers, strings.TrimSpace(th.Text()))
		}
	})

	// Find the first <tr> in <tbody> (the 'Spaces' row)
	spacesRow := doc.Find("tbody tr").First()
	spacesRow.Find("td").Each(func(i int, td *goquery.Selection) {
		if i == 0 {
			return // skip the first column ("Spaces")
		}
		// Check if this <td> contains <i class="fa-solid fa-circle"></i>
		if td.Find("i.fa-solid.fa-circle").Length() > 0 {
			if i-1 < len(headers) {
				regions = append(regions, headers[i-1])
			}
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
				if len(regionCode) > 0 {
					regionMap[regionCode] = regionName
				}
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
