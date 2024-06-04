package handler

import (
	"net/http"

	"github.com/sb-nour/providers-endpoints/lib"
	gee "github.com/tbxark/g4vercel"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	server := gee.New()
	server.GET("/", func(context *gee.Context) {
		regions := lib.GetRegions()
		context.JSON(200, regions)
	})
	server.Handle(w, r)
}
