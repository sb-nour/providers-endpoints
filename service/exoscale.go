package service

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func getExoscaleStorageRegions() map[string]string {
	url := "https://community.exoscale.com/documentation/platform/exoscale-datacenter-zones/"

	// Make a GET request to the URL
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making GET request:", err)
		return nil
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		fmt.Println("Error loading HTML:", err)
		return nil
	}

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
	return Regions{
		Storage: getExoscaleStorageRegions(),
		Compute: getExoscaleStorageRegions(),
	}
}
