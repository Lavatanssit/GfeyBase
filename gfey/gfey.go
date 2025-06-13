package gfey

import (
	"net/http"
)

// HandlerFunc 定义所有的请求处理函数 参数是 writer 和 *request
type HandlerFunc func(*Context)

// Engine 实现了 ServeHTTP 接口
type Engine struct {
	router *Router
}

// New 是 Engine 构造函数
func New() *Engine {
	return &Engine{router: NewRouter()}
}

// AddRoute 用于向路由表中添加路由
func (engine *Engine) AddRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.AddRoute(method, pattern, handler)
}

// GET 用于添加 HTTP 的 GET 请求
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.AddRoute("GET", pattern, handler)
}

// POST 用于添加 HTTP 的 GET 请求
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.AddRoute("POST", pattern, handler)
}

// Run 用于运行 http 服务器，调用标准库 ListenAndServe 方法
func (engine *Engine) Run(address string) (err error) {
	return http.ListenAndServe(address, engine)
}

// Engine 实现标准库的 ServeHTTP 接口
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := NewContext(w, req)
	engine.router.Handle(c)
}
