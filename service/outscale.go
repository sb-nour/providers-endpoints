package service

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getOutscaleStorageRegions(doc *goquery.Document) map[string]string {
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
func getOutscaleComputeRegions(doc *goquery.Document) map[string]string {

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
		regionMap[regionCode] = fmt.Sprintf("%s - %s", regionName, regionCode)
	})

	return regionMap
}

func GetOutscaleRegions() Regions {
	doc, _ := get("https://docs.outscale.com/en/userguide/Regions-Endpoints-and-Subregions-Reference.html")
	return Regions{
		Storage: getOutscaleStorageRegions(doc),
		Compute: getOutscaleComputeRegions(doc),
	}
}
