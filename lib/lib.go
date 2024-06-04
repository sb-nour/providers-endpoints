package lib

import (
	"sync"

	"github.com/sb-nour/providers-endpoints/service"
)

var Providers = []struct {
	name string
	fn   func() service.Regions
}{
	{"Amazon AWS", service.GetAmazonRegions},
	{"Amazon Lightsail", service.GetLightsailRegions},
	{"DigitalOcean", service.GetDigitalOceanRegions},
	{"UpCloud", service.GetUpcloudRegions},
	{"Exoscale", service.GetExoscaleRegions},
	// {"Wasabi", service.GetWasabiRegions},
	{"Google Cloud", service.GetGoogleCloudRegions},
	{"Backblaze", service.GetBackblazeRegions},
	{"Linode", service.GetLinodeRegions},
	{"Outscale", service.GetOutscaleRegions},
	{"Storj", service.GetStorjRegions},
	{"Vultr", service.GetVultrRegions},
	{"Hetzner", service.GetHetznerRegions},
	{"Synology", service.GetSynologyRegions},
}

func GetRegions() map[string]service.Regions {
	workerCount := 10
	regions := make(map[string]service.Regions)
	var wg sync.WaitGroup
	providerRegions := make(chan service.ProviderRegions, len(Providers))
	workerPool := make(chan struct{}, workerCount)

	for _, provider := range Providers {
		workerPool <- struct{}{}
		wg.Add(1)
		go func(provider struct {
			name string
			fn   func() service.Regions
		}) {
			defer func() {
				<-workerPool
				wg.Done()
			}()
			providerRegions <- service.ProviderRegions{Provider: provider.name, Regions: provider.fn()}
		}(provider)
	}

	go func() {
		wg.Wait()
		close(providerRegions)
	}()

	for pr := range providerRegions {
		regions[pr.Provider] = pr.Regions
	}

	return regions
}
