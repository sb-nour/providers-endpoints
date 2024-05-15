package service

import (
	"fmt"
	"strings"
	"sync"

	"github.com/go-ping/ping"
)

var synologyKnownPrefixes = []string{
	"us-",
	"eu-",
	"ap-",
	"ca-",
	"sa-",
}

// operation in ['create-snapshot', 'delete-snapshot', 'get-snapshot', 'list-snapshots', 'list-instances', 'get-instance']
func transformLabelSynology(regionCode string) string {
	parts := strings.Split(regionCode, "-")
	return fmt.Sprintf("%s %s - %s", strings.ToUpper(parts[0]), strings.Title(parts[1]), regionCode)
}

func getSynologyStorageRegions() map[string]string {

	var wg sync.WaitGroup
	results := make(chan map[string]string)

	for _, region := range synologyKnownPrefixes {
		for i := 0; i < 6; i++ {
			wg.Add(1)
			go func(region string, i int) {
				defer wg.Done()
				formatedRegionCode := region + fmt.Sprintf("%03d", i)
				endpoint := formatedRegionCode + ".s3.synologyc2.net"
				pinger, err := ping.NewPinger(endpoint)
				if err != nil {
					return
				}
				pinger.Count = 1
				pinger.Run() // blocks until finished
				results <- map[string]string{formatedRegionCode: transformLabelSynology(formatedRegionCode)}
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

func GetSynologyRegions() Regions {
	return Regions{
		Storage: getSynologyStorageRegions(),
	}
}
