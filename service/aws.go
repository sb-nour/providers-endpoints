package service

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getAmazonS3Regions() (map[string]string, error) {
	url := "https://docs.aws.amazon.com/general/latest/gr/s3.html"
	doc, err := get(url)
	if err != nil {
		return nil, err
	}

	var regionMap map[string]string = make(map[string]string)

	doc.Find("#main-col-body div table").Each(func(i int, table *goquery.Selection) {
		// CHeck if the thead first row's first th is "Region Name"
		if table.Find("thead th").Eq(0).Text() != "Region Name" {
			return
		}
		table.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
			if row.Children().Length() == 5 {
				regionCode := strings.Trim(row.Children().Eq(1).Text(), " \n")
				regionName := strings.Trim(row.Children().Eq(0).Text(), " \n") + " - " + strings.Trim(row.Children().Eq(1).Text(), " \n")

				regionMap[regionCode] = regionName
			}
		})
	})

	return regionMap, nil
}

func getAmazonEC2Regions() (map[string]string, error) {

	url := "https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html#concepts-regions"
	doc, err := get(url)
	if err != nil {
		return nil, err
	}

	var regionMap map[string]string = make(map[string]string)

	doc.Find("table").Each(func(i int, table *goquery.Selection) {
		if table.Find("thead th").Length() == 3 && table.Find("tbody tr").Length() > 5 {
			table.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
				regionCode := row.Find("td").Eq(0).Text()
				regionName := row.Find("td").Eq(1).Text()
				regionMap[regionCode] = fmt.Sprintf("%s - %s", regionName, regionCode)
			})
			return
		}
	})

	return regionMap, nil
}

func GetAmazonRegions() Regions {

	s3Regions, err := getAmazonS3Regions()
	if err != nil {
		// Throw error
		fmt.Println(err)
	}
	ec2Regions, err := getAmazonEC2Regions()
	if err != nil {
		// Throw error
		fmt.Println(err)
	}

	regions := Regions{
		Storage: s3Regions,
		Compute: ec2Regions,
	}

	return regions
}
