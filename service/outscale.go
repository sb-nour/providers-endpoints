package service

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getOutscaleStorageRegions(doc *goquery.Document) map[string]string {
	regionMap := make(map[string]string)
	var currentRegion string
	doc.Find("h2#_mapping_between_subregions_and_physical_zones").NextAllFiltered("div.sectionbody").First().
		Find("table.tableblock tbody tr").Each(func(i int, row *goquery.Selection) {
		cols := row.Find("td")
		if cols.Length() == 3 {
			regionCell := cols.Eq(0).Find("p").Text()
			if regionCell != "" {
				currentRegion = strings.TrimSpace(regionCell)
			}
			subregions := strings.TrimSpace(cols.Eq(1).Find("p").Text())
			physicalZone := strings.TrimSpace(cols.Eq(2).Find("p").Text())
			regionMap[currentRegion] = fmt.Sprintf("Region: %s - Subregions: %s - Physical Zones: %s", currentRegion, subregions, physicalZone)
		}
	})
	return regionMap
}

func getOutscaleComputeRegions(doc *goquery.Document) map[string]string {
	// Same logic as storage, as the table now contains all region info
	return getOutscaleStorageRegions(doc)
}

func GetOutscaleRegions() Regions {
	doc, _ := get("https://docs.outscale.com/en/userguide/About-Regions-and-Subregions.html")
	return Regions{
		Storage: getOutscaleStorageRegions(doc),
		Compute: getOutscaleComputeRegions(doc),
	}
}
