package main

import "regexp"

// regionMapping is a map that stores the translation of region codes to their corresponding names.
var regionMapping = map[string]string{
	"NYC": "New York City",
	"AMS": "Amsterdam",
	"SFO": "San Francisco",
	"SGP": "Singapore",
	"FRA": "Frankfurt",
	"BLR": "Bangalore",
	"SYD": "Sydney",
}

// translateRegionCode translates a region code into a city name.
// It removes any numbers from the code and looks up the city name using the cleaned code.
// If the city name is found, it is returned. Otherwise, "Region code not found" is returned.
func translateRegionCode(code string) string {
	// Use regular expression to remove numbers
	re := regexp.MustCompile("[0-9]+")
	cleanCode := re.ReplaceAllString(code, "")

	// Lookup the city name using the cleaned code
	city, exists := regionMapping[cleanCode]
	if exists {
		return city
	}
	return "Region code not found"
}

// translateRegions translates a slice of region codes into a key-value pair of region codes and city names.
// It iterates over the regions and calls translateRegionCode to get the city name for each region.
func translateRegions(regions []string) map[string]string {
	translatedRegions := make(map[string]string)

	for _, region := range regions {
		city := translateRegionCode(region)
		translatedRegions[region] = city
	}

	return translatedRegions
}
