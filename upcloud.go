package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getUpcloudStorageRegions() map[string]string {
	url := "https://upcloud.com/data-centres"

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

func getUpcloudRegions() Regions {
	return Regions{
		Storage: getUpcloudStorageRegions(),
	}
}
