package service

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getHetznerRegions() map[string]string {
	url := "https://docs.hetzner.com/cloud/general/locations/"
	doc, _ := get(url)

	regionMap := make(map[string]string)

	// Select the first table in the document
	table := doc.Find("table").First()
	table.Find("tbody tr").Each(func(i int, tr *goquery.Selection) {
		tr.Find("td").Each(func(j int, td *goquery.Selection) {
			// Skip empty cells
			if strings.TrimSpace(td.Text()) == "" {
				return
			}

			codeSel := td.Find("code")
			if codeSel.Length() > 0 {
				regionCode := strings.TrimSpace(codeSel.Text())
				// Extract the location name by removing the code part
				fullText := strings.TrimSpace(td.Text())
				locationName := strings.TrimSpace(strings.Replace(fullText, regionCode, "", 1))

				// Clean up the location name
				if locationName != "" {
					regionMap[regionCode] = fmt.Sprintf("%s - %s", locationName, regionCode)
				}
			}
		})
	})

	return regionMap
}

// Helper to strip HTML tags from a string
func stripTags(html string) string {
	var result strings.Builder
	inTag := false
	for _, r := range html {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			continue
		}
		if !inTag {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func GetHetznerRegions() Regions {
	regions := getHetznerRegions()
	return Regions{
		Compute: regions,
		Storage: regions,
	}
}
