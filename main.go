package main

import (
	"gfey"
	"net/http"
)

func main() {
	r := gfey.New()

	r.GET("/", func(ctx *gfey.Context) {
		ctx.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/news", func(ctx *gfey.Context) {
			// url example: /news?title=TimesNews&content=ConferencesStart
			ctx.String(http.StatusOK, "title:%s\ncontent:%s\n", ctx.Query("title"), ctx.Query("content"))
		})
		v1.POST("/login", func(ctx *gfey.Context) {
			ctx.Json(http.StatusOK, gfey.H{
				"username": ctx.PostForm("username"),
				"password": ctx.PostForm("password"),
			})
		})
	}

	v2 := r.Group("/v2")
	{
		v2.GET("/news/:date", func(ctx *gfey.Context) {
			ctx.String(http.StatusOK, "news in date %v\n", ctx.Param("date"))
		})

		v2.GET("/news/:date/*newsName", func(ctx *gfey.Context) {
			ctx.String(http.StatusOK, "news in date %v\n", ctx.Param("date"))
			ctx.Json(http.StatusOK, gfey.H{
				"date":     ctx.Param("date"),
				"newsName": ctx.Param("newsName"),
			})
		})
	}

	r.Run(":9999")
}
