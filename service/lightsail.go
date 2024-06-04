package service

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getLightsailComputeRegions() map[string]string {
	url := "https://docs.aws.amazon.com/lightsail/latest/userguide/understanding-regions-and-availability-zones-in-amazon-lightsail.html"
	doc, _ := get(url)

	var regionMap map[string]string = make(map[string]string)

	doc.Find(".listitem").Each(func(i int, listItem *goquery.Selection) {
		value := listItem.Find("p").Text()

		// there are two pairs of parentheses in the value, extract the regionCode from the second pair
		regionCode := value[strings.LastIndex(value, "(")+1 : strings.LastIndex(value, ")")]
		regionName := value

		regionMap[regionCode] = regionName
	})

	return regionMap
}

func GetLightsailRegions() Regions {
	return Regions{
		Compute: getLightsailComputeRegions(),
	}
}
