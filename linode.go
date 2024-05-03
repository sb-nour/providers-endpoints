package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getLinodeStorageRegions() map[string]string {
	url := "https://www.linode.com/docs/products/storage/object-storage/"

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
		if table.Find("tbody tr").Length() < 5 {
			return
		}

		if table.Find("thead th").Length() == 2 && table.Find("thead th").Eq(0).Text() == "Data Center" {

			table.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
				regionCode := row.Find("td").Eq(1).Text()
				regionName := strings.ReplaceAll(row.Find("td").Eq(0).Text(), "*", "")
				regionMap[regionCode] = regionName
			})
		}
	})

	return regionMap
}

func getLinodeRegions() Regions {
	return Regions{
		Storage: getLinodeStorageRegions(),
	}
}
