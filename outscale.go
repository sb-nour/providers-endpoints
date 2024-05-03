package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getOutscaleStorageRegions() map[string]string {
	url := "https://docs.outscale.com/en/userguide/Regions-Endpoints-and-Subregions-Reference.html"

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
	doc.Find("#_outscale_object_storage_oos").First().Next().Find("tbody tr").Each(func(i int, row *goquery.Selection) {

		regionCode := strings.ToLower(row.Find("td").Eq(0).Text())
		// parts := strings.Split(row.Find("td").Eq(0).Text(), " - ")
		regionName := row.Find("td").Eq(0).Text()
		parts := strings.Split(regionName, "-")
		// regionName = fmt.Sprintf("%s %s %s", strings.ToUpper(parts[0]), strings.ToTitle(parts[1]), parts[2])

		if len(parts) == 3 {
			regionName = fmt.Sprintf("%s %s %s", strings.ToUpper(parts[0]), strings.Title(parts[1]), parts[2])
		} else {
			regionName = fmt.Sprintf("%s %s %s %s", strings.Title(parts[0]), strings.ToTitle(parts[1]), strings.ToUpper(parts[2]), parts[3])
		}
		regionMap[regionCode] = regionName
	})

	return regionMap
}
func getOutscaleComputeRegions() map[string]string {
	url := "https://docs.outscale.com/en/userguide/Regions-Endpoints-and-Subregions-Reference.html"

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
	doc.Find("#_available_endpoints").First().Next().Find("tbody tr").Each(func(i int, row *goquery.Selection) {

		regionCode := strings.ToLower(row.Find("td").Eq(0).Text())
		// parts := strings.Split(row.Find("td").Eq(0).Text(), " - ")
		regionName := row.Find("td").Eq(0).Text()
		parts := strings.Split(regionName, "-")
		// regionName = fmt.Sprintf("%s %s %s", strings.ToUpper(parts[0]), strings.ToTitle(parts[1]), parts[2])

		if len(parts) == 3 {
			regionName = fmt.Sprintf("%s %s %s", strings.ToUpper(parts[0]), strings.Title(parts[1]), parts[2])
		} else {
			regionName = fmt.Sprintf("%s %s %s %s", strings.Title(parts[0]), strings.ToTitle(parts[1]), strings.ToUpper(parts[2]), parts[3])
		}
		regionMap[regionCode] = regionName
	})

	return regionMap
}

func getOutscaleRegions() Regions {
	return Regions{
		Storage: getOutscaleStorageRegions(),
		Compute: getOutscaleComputeRegions(),
	}
}
