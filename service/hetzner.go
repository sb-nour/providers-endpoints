package service

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getHetznerRegions() map[string]string {
	url := "https://docs.hetzner.com/cloud/general/locations/"
	doc, _ := get(url)

	var regionMap map[string]string = make(map[string]string)

	doc.Find("table tbody tr").Each(func(i int, tr *goquery.Selection) {
		if tr.Children().Length() != 3 {
			return
		}

		tr.Find("td").Each(func(i int, td *goquery.Selection) {
			if td.Find("code") != nil && td.Text() != "" {
				regionCode := td.Find("code").Text()
				regionName := strings.TrimSpace(td.Contents().Not("code, br").Text())
				regionMap[regionCode] = fmt.Sprintf("%s - %s", regionName, regionCode)
			}
		})
	})

	return regionMap
}

func GetHetznerRegions() Regions {
	regions := getHetznerRegions()
	return Regions{
		Compute: regions,
		Storage: regions,
	}
}
