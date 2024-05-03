package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/go-ping/ping"
)

var knownPrefixes = []string{
	"us-east-",
	"us-west-",
	"eu-central-",
	"eu-west-",
	"ap-south-",
	"ap-northeast-",
	"ap-southeast-",
	"ca-central-",
	"sa-east-",
}

// operation in ['create-snapshot', 'delete-snapshot', 'get-snapshot', 'list-snapshots', 'list-instances', 'get-instance']
func transformLabel(regionCode string) string {
	parts := strings.Split(regionCode, "-")
	thirdPart, _ := strconv.Atoi(parts[2])
	return fmt.Sprintf("%s %s %1d - %s", strings.ToUpper(parts[0]), strings.Title(parts[1]), thirdPart, regionCode)
}

func getBackblazeStorageRegions() map[string]string {

	var wg sync.WaitGroup
	results := make(chan map[string]string)

	for _, region := range knownPrefixes {
		for i := 0; i < 6; i++ {
			wg.Add(1)
			go func(region string, i int) {
				defer wg.Done()
				formatedRegionCode := region + fmt.Sprintf("%03d", i)
				endpoint := "s3." + formatedRegionCode + ".backblazeb2.com"
				pinger, err := ping.NewPinger(endpoint)
				if err != nil {
					return
				}
				pinger.Count = 1
				pinger.Run() // blocks until finished
				results <- map[string]string{formatedRegionCode: transformLabel(formatedRegionCode)}
			}(region, i)
		}
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	regionMap := make(map[string]string)
	for result := range results {
		for k, v := range result {
			regionMap[k] = v
		}
	}

	return regionMap
}

func getBackblazeRegions() Regions {
	return Regions{
		Storage: getBackblazeStorageRegions(),
	}
}
