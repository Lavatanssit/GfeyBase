package middlewares

import (
	"gfey"
	"log"
	"time"
)

// Logger: 日志中间件，用于的打印请求处理的时间
func Logger() gfey.HandlerFunc {
	return func(ctx *gfey.Context) {
		// 起始时间
		start_t := time.Now()
		// 链式执行中间件、处理请求
		ctx.Next()
		// 计算处理时间
		time.Sleep(100 * time.Millisecond)
		log.Printf("[%d] %s in %v", ctx.StatusCode, ctx.Req.RequestURI, time.Since(start_t))
	}
}
