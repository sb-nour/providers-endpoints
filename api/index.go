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
	{"AWS", service.GetAmazonRegions},
	{"BACKBLAZE", service.GetBackblazeRegions},
	{"DIGITALOCEAN", service.GetDigitalOceanRegions},
	{"EXOSCALE", service.GetExoscaleRegions},
	{"GOOGLE_CLOUD", service.GetGoogleCloudRegions},
	{"LIGHTSAIL", service.GetLightsailRegions},
	{"LINODE", service.GetLinodeRegions},
	{"OUTSCALE", service.GetOutscaleRegions},
	{"STORJ", service.GetStorjRegions},
	{"UPCLOUD", service.GetUpcloudRegions},
	{"VULTR", service.GetVultrRegions},
	{"WASABI", service.GetWasabiRegions},
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
