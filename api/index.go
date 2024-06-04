package handler

import (
	"net/http"
	"strings"

	gee "github.com/tbxark/g4vercel"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	server := gee.New()
	server.GET("/", func(context *gee.Context) {
		regions := GetRegions()
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
