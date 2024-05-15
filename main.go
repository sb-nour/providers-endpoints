package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/sb-nour/providers-endpoints/service"
)

type ProviderRegions struct {
	Provider string
	Regions  service.Regions
}

func getRegions() map[string]service.Regions {
	var wg sync.WaitGroup
	providerRegions := make(chan ProviderRegions)

	providers := []struct {
		name string
		fn   func() service.Regions
	}{
		{"Amazon AWS", service.GetAmazonRegions},
		{"Amazon Lightsail", service.GetLightsailRegions},
		{"DigitalOcean", service.GetDigitalOceanRegions},
		{"UpCloud", service.GetUpcloudRegions},
		{"Exoscale", service.GetExoscaleRegions},
		{"Wasabi", service.GetWasabiRegions},
		{"Google Cloud", service.GetGoogleCloudRegions},
		{"Backblaze", service.GetBackblazeRegions},
		{"Linode", service.GetLinodeRegions},
		{"Outscale", service.GetOutscaleRegions},
		{"Storj", service.GetStorjRegions},
		{"Vultr", service.GetVultrRegions},
		{"Hetzner", service.GetHetznerRegions},
		{"Synology", service.GetSynologyRegions},
	}

	for _, provider := range providers {
		wg.Add(1)
		go func(provider struct {
			name string
			fn   func() service.Regions
		}) {
			defer wg.Done()
			providerRegions <- ProviderRegions{Provider: provider.name, Regions: provider.fn()}
		}(provider)
	}

	go func() {
		wg.Wait()
		close(providerRegions)
	}()

	regions := make(map[string]service.Regions)
	for pr := range providerRegions {
		regions[pr.Provider] = pr.Regions
	}

	return regions
}

func main() {
	regions := getRegions()
	regionsJson, err := json.Marshal(regions)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	fmt.Println(string(regionsJson))
}
