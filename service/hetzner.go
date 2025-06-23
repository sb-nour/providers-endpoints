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
			codeSel := td.Find("code")
			if codeSel.Length() > 0 {
				regionCode := strings.TrimSpace(codeSel.Text())
				// Remove the <code>...</code> from the HTML to get the location name
				locationHtml, _ := td.Html()
				// Remove the code tag and trim
				locationName := strings.TrimSpace(strings.Replace(td.Text(), regionCode, "", 1))
				// If locationName is still empty, fallback to the text before <code>
				if locationName == "" {
					parts := strings.Split(locationHtml, "<code>")
					if len(parts) > 0 {
						locationName = strings.TrimSpace(stripTags(parts[0]))
					}
				}
				regionMap[regionCode] = fmt.Sprintf("%s - %s", locationName, regionCode)
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
