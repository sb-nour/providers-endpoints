package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getAmazonS3Regions() map[string]string {
	url := "https://docs.aws.amazon.com/general/latest/gr/s3.html"

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

	doc.Find("#main-col-body div table tbody tr").Each(func(i int, row *goquery.Selection) {
		// Check if the row has five columns
		if row.Children().Length() == 5 {

			// First Column is Region Name,
			// Second Column is Region Code

			// create a map of region code to region name

			regionCode := strings.Trim(row.Children().Eq(1).Text(), " \n")
			regionName := strings.Trim(row.Children().Eq(0).Text(), " \n") + " - " + strings.Trim(row.Children().Eq(1).Text(), " \n")

			regionMap[regionCode] = regionName
		}
	})

	return regionMap
}

func getAmazonEC2Regions() map[string]string {

	url := "https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html#concepts-regions"

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
		if table.Find("thead th").Length() == 3 && table.Find("tbody tr").Length() > 5 {
			table.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
				regionCode := row.Find("td").Eq(0).Text()
				regionName := row.Find("td").Eq(1).Text()
				regionMap[regionCode] = regionName
			})
			return
		}
	})

	return regionMap
}

func getAmazonRegions() Regions {
	// Get Amazon S3 Regions
	s3Regions := getAmazonS3Regions()

	// Get Amazon EC2 Regions
	ec2Regions := getAmazonEC2Regions()

	regions := Regions{
		Storage: s3Regions,
		Compute: ec2Regions,
	}

	return regions
}
