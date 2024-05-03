package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getGoogleCloudStorageRegions() map[string]string {
	url := "https://cloud.google.com/storage/docs/locations/"

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

func getGoogleCloudRegions() Regions {
	return Regions{
		Storage: getGoogleCloudStorageRegions(),
	}
}
