package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getUpcloudStorageRegions(doc *goquery.Document) map[string]string {
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
func getUpcloudComputeRegions(doc *goquery.Document) map[string]string {

	var regionMap map[string]string = make(map[string]string)

	doc.Find(".accordion").First().Find(".accordion-item").Each(func(i int, item *goquery.Selection) {
		// Check if item has an <li> tag with "Object Storage" text
		if strings.Contains(item.Find("li").Text(), "Cloud Servers") {

			regionCode := strings.ToLower(item.Find("button h3").First().Text())
			regionSplit := strings.Split(item.Find("button .location").First().Text(), ", ")
			regionName := regionSplit[1] + " - " + regionSplit[0] + " - " + regionCode

			regionMap[regionCode] = regionName
		}
	})

	return regionMap
}

func GetUpcloudRegions() Regions {
	doc, err := get("https://upcloud.com/data-centres")
	if err != nil {
		// fmt.Printf("Error: %v", err)
		// Load the regions from the local file ./upcloud_fallback.json
		filePath := "./service/upcloud_fallback.json"
		absPath, _ := filepath.Abs(filePath)
		jsonContent, err := ioutil.ReadFile(absPath)
		if err != nil {
			if debugging {
				fmt.Printf("Error reading file: %v", err)
			}
			return Regions{}
		}

		// json has "storage" and "compute" keys
		var regions map[string]map[string]string
		json.Unmarshal(jsonContent, &regions)

		return Regions{
			Storage: regions["storage"],
			Compute: regions["compute"],
		}
	}
	storageRegions := getUpcloudStorageRegions(doc)
	computeRegions := getUpcloudComputeRegions(doc)

	return Regions{
		Storage: storageRegions,
		Compute: computeRegions,
	}
}
