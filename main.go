package main

import (
	"gfey"
	"net/http"
)

func main() {
	r := gfey.New()

	r.GET("/", func(ctx *gfey.Context) {
		ctx.HTML(http.StatusOK, "<h1>index</h1>")
	})

	r.GET("/news", func(ctx *gfey.Context) {
		// url example: /news?title=TimesNews&content=ConferencesStart
		ctx.String(http.StatusOK, "title:%s\ncontent:%s\n", ctx.Query("title"), ctx.Query("content"))
	})

	r.POST("/login", func(ctx *gfey.Context) {
		ctx.Json(http.StatusOK, gfey.H{
			"username": ctx.PostForm("username"),
			"password": ctx.PostForm("password"),
		})
	})

	r.Run(":9999")
}
