package service

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Transforms 'ap-northeast-1' => 'AP Northeast 1 (Tokyo) - ap-northeast-1'
// Tokyo is the region name and ap-northeast-1 is the region code
func transformRegionName(regionName string, regionCode string) string {
	// transform ap-northeast-1 to AP Northeast 1
	splitRegionName := strings.Split(strings.Title(strings.ReplaceAll(regionCode, "-", " ")), " ")

	finalSplitRegionName := make([]string, 0)
	for i, region := range splitRegionName {
		if i == 0 {
			finalSplitRegionName = append(finalSplitRegionName, strings.ToUpper(region))
			continue
		}
		finalSplitRegionName = append(finalSplitRegionName, strings.Title(region))
	}
	// append the region code to the region name
	return fmt.Sprintf("%s (%s) - %s", strings.Join(finalSplitRegionName, " "), regionName, regionCode)
}

// getWasabiRegions retrieves the regions and their corresponding codes from the Wasabi website.
// It makes a GET request to the Wasabi locations page and parses the HTML response to extract the region information.
// The regions and their codes are stored in a map[string]string, where the region code is the key and the region name is the value.
// If an error occurs during the HTTP request or HTML parsing, nil is returned.
func getWasabiStorageRegions() map[string]string {
	url := "https://wasabi.com/company/storage-regions"
	doc, _ := get(url)

	if debugging {
		fmt.Println("Wasabi storage regions URL: ", url)
	}

	var regionMap map[string]string = make(map[string]string)

	// Iterate over each row in the table body
	doc.Find("tbody .c-table-row").Each(func(index int, row *goquery.Selection) {
		if debugging {
			fmt.Println("Row: ", index)
		}
		// For each row, find the cells with class 'c-table-cell'
		row.Find(".c-table-cell").Each(func(cellIndex int, cell *goquery.Selection) {
			// Extract the strong text as the region name
			regionName := cell.Find("strong").Text()
			// Extract the region code
			regionCode := strings.TrimSpace(cell.Find("div").Contents().Not("strong, br").Text())
			// if "region" exists in the region code, remove it
			regionCode = strings.TrimSpace(strings.Replace(regionCode, "region", "", -1))
			// if an "&" exists in the region code split it and add both parts
			if strings.Contains(regionCode, "&") {
				regions := strings.Split(regionCode, " & ")
				for _, region := range regions {
					regionMap[region] = transformRegionName(regionName, region)
				}
				return
			}
			// Assign the region code to the region name in the map
			// Add the region code and name to the map
			regionMap[regionCode] = transformRegionName(regionName, regionCode)
		})
	})

	return regionMap
}

func GetWasabiRegions() Regions {
	return Regions{
		Storage: getWasabiStorageRegions(),
	}
}
