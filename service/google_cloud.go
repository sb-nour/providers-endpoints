package service

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getGoogleCloudStorageRegions() map[string]string {
	url := "https://cloud.google.com/storage/docs/locations/"
	doc, _ := get(url)

	var regionMap map[string]string = make(map[string]string)
	doc.Find("table").Each(func(i int, table *goquery.Selection) {
		// if table doesn't have more than 2 rows, return
		if table.Find("tbody tr").Length() < 2 {
			return
		}

		if table.Find("thead th").Length() != 4 {
			return
		}
		currentRegion := ""
		table.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
			// check if td:nth-child(1) is not empty
			if row.Find("td").Eq(1).Text() != "" {
				regionCode := strings.ToLower(row.Find("td").Eq(1).Text())
				// regionName := row.Find("td").Eq(2).Text()
				key := currentRegion + " - " + regionCode
				regionMap[regionCode] = key
			} else {
				currentRegion = row.Find("td").Eq(0).Text()
			}
		})

	})

	return regionMap
}
func getGoogleCloudComputeRegions() map[string]string {
	url := "https://cloud.google.com/compute/docs/regions-zones"
	doc, _ := get(url)

	var regionMap map[string]string = make(map[string]string)
	doc.Find("table").Each(func(i int, table *goquery.Selection) {
		if table.Find("thead th").Length() == 6 {
			table.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
				regionCode := strings.ToLower(row.Find("td").Eq(0).Text())
				regionName := row.Find("td").Eq(1).Text()
				regionMap[regionCode] = fmt.Sprintf("%s - %s", regionName, regionCode)
			})
		}

	})

	return regionMap
}

func GetGoogleCloudRegions() Regions {
	return Regions{
		Storage: getGoogleCloudStorageRegions(),
		Compute: getGoogleCloudComputeRegions(),
	}
}
