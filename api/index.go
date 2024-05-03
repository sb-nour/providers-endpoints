package handler

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/sb-nour/providers-endpoints/service"
	gee "github.com/tbxark/g4vercel"
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
		{"AWS", service.GetAmazonRegions},
		{"LIGHTSAIL", service.GetLightsailRegions},
		{"DIGITALOCEAN", service.GetDigitalOceanRegions},
		{"UPCLOUD", service.GetUpcloudRegions},
		{"EXOSCALE", service.GetExoscaleRegions},
		{"WASABI", service.GetWasabiRegions},
		{"GOOGLE_CLOUD", service.GetGoogleCloudRegions},
		{"BACKBLAZE", service.GetBackblazeRegions},
		{"LINODE", service.GetLinodeRegions},
		{"OUTSCALE", service.GetOutscaleRegions},
		{"STORJ", service.GetStorjRegions},
		{"VULTR", service.GetVultrRegions},
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

func Handler(w http.ResponseWriter, r *http.Request) {
	server := gee.New()

	server.GET("/", func(context *gee.Context) {
		regions := getRegions()
		context.JSON(200, regions)
	})
	server.GET("/hello", func(context *gee.Context) {
		name := context.Query("name")
		if name == "" {
			context.JSON(400, gee.H{
				"message": "name not found",
			})
		} else {
			context.JSON(200, gee.H{
				"data": fmt.Sprintf("Hello %s!", name),
			})
		}
	})
	server.GET("/user/:id", func(context *gee.Context) {
		context.JSON(400, gee.H{
			"data": gee.H{
				"id": context.Param("id"),
			},
		})
	})
	server.GET("/long/long/long/path/*test", func(context *gee.Context) {
		context.JSON(200, gee.H{
			"data": gee.H{
				"url": context.Path,
			},
		})
	})
	server.Handle(w, r)
}
