package main

import (
	"gfey"
	"log"
	"net/http"
	"time"
)

func onlyForV2() gfey.HandlerFunc {
	return func(ctx *gfey.Context) {
		// 起始时间
		start_t := time.Now()
		// 链式执行中间件、处理请求
		ctx.Next()
		// 计算处理时间
		log.Printf("Group V2 ：[%d] %s in %v", ctx.StatusCode, ctx.Req.RequestURI, time.Since(start_t).Milliseconds())
	}
}

func main() {
	r := gfey.New()

	r.Use(gfey.Logger())

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
	v2.Use(onlyForV2())
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
