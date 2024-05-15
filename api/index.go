package handler

import (
	"net/http"
	"strings"
	"sync"

	"github.com/sb-nour/providers-endpoints/service"
	gee "github.com/tbxark/g4vercel"
)

var providers = []struct {
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

func getRegions() map[string]service.Regions {
	var wg sync.WaitGroup
	providerRegions := make(chan service.ProviderRegions)

	for _, provider := range providers {
		wg.Add(1)
		go func(provider struct {
			name string
			fn   func() service.Regions
		}) {
			defer wg.Done()
			providerRegions <- service.ProviderRegions{Provider: provider.name, Regions: provider.fn()}
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

func Handler(w http.ResponseWriter, r *http.Request) {
	server := gee.New()
	server.GET("/", func(context *gee.Context) {
		regions := getRegions()
		context.JSON(200, regions)
	})
	server.GET("/:key", func(context *gee.Context) {
		key := strings.ToUpper(context.Param("key"))
		// if key is in `providers`, run the function and return the result
		for _, provider := range providers {
			if key == provider.name {
				context.JSON(200, provider.fn())
				break
			}
		}
	})
	server.Handle(w, r)
}
