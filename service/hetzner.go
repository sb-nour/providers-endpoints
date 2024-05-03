package service

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getHetznerRegions() map[string]string {
	url := "https://docs.hetzner.com/cloud/general/locations/"

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
	return Regions{
		Compute: getHetznerRegions(),
		Storage: getHetznerRegions(),
	}
}
