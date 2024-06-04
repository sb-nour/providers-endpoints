package service

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

func getExoscaleStorageRegions() map[string]string {
	url := "https://community.exoscale.com/documentation/platform/exoscale-datacenter-zones/"
	doc, _ := get(url)

	var regionMap map[string]string = make(map[string]string)
	doc.Find("table").Each(func(i int, table *goquery.Selection) {
		// if table doesn't have more than 2 rows, return
		if table.Find("tbody tr").Length() < 2 {
			return
		}

		table.Find("tbody tr").Each(func(i int, row *goquery.Selection) {

			regionCode := row.Find("td").Eq(2).Text()
			regionName := row.Find("td").Eq(1).Text()
			regionMap[regionCode] = fmt.Sprintf("%s - %s", regionName, regionCode)
		})

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
