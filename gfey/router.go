package gfey

import (
	"log"
	"net/http"
)

// Router: 路由类，存储 url 到处理句柄的映射
type Router struct {
	handlers map[string]HandlerFunc
}

// NewRouter: 构造Router，返回指向新Router实例的指针
func NewRouter() *Router {
	return &Router{make(map[string]HandlerFunc)}
}

// AddRoute: 向router中添加路由项
func (r *Router) AddRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	r.handlers[key] = handler
}

// Handle: 在路由表中查找上下文请求，并执行。查找不到返回404
func (r *Router) Handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %v\n", c.Path)
	}
}
