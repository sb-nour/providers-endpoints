package handler

import (
	"fmt"
	"net/http"

	gee "github.com/tbxark/g4vercel"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	server := gee.New()

	server.GET("/", func(context *gee.Context) {
		context.JSON(200, gee.H{
			"message": "hello go from vercel !!!!",
		})
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