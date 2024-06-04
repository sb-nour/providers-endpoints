package service

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
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
	iterations := 8
	workerCount := len(knownPrefixes) * iterations // Set the number of concurrent workers

	// Create a buffered channel to limit the number of concurrent goroutines
	workerPool := make(chan struct{}, workerCount)

	for _, region := range knownPrefixes {
		for i := 0; i < iterations; i++ {
			wg.Add(1)
			workerPool <- struct{}{} // Acquire a worker slot
			go func(region string, i int) {
				defer func() {
					wg.Done()
					<-workerPool // Release the worker slot
				}()

				formatedRegionCode := region + fmt.Sprintf("%03d", i)
				endpoint := "s3." + formatedRegionCode + ".backblazeb2.com"
				addrs, err := net.LookupHost(endpoint)
				if err != nil {
					return
				}
				if len(addrs) == 0 {
					return
				}

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

func GetBackblazeRegions() Regions {
	return Regions{
		Storage: getBackblazeStorageRegions(),
	}
}
