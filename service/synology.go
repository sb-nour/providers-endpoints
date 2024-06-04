package service

import (
	"fmt"
	"net"
	"strings"
	"sync"
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
	iterations := 8
	workerCount := len(synologyKnownPrefixes) * iterations // Set the number of concurrent workers

	// Create a buffered channel to limit the number of concurrent goroutines
	workerPool := make(chan struct{}, workerCount)

	for _, region := range synologyKnownPrefixes {
		for i := 0; i < iterations; i++ {
			wg.Add(1)
			workerPool <- struct{}{} // Acquire a worker slot

			go func(region string, i int) {
				defer func() {
					wg.Done()
					<-workerPool // Release the worker slot
				}()

				formatedRegionCode := region + fmt.Sprintf("%03d", i)
				endpoint := formatedRegionCode + ".s3.synologyc2.net"
				addrs, err := net.LookupHost(endpoint)
				if err != nil {
					return
				}
				if len(addrs) == 0 {
					return
				}

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
