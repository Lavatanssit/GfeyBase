package gfey

import (
	"log"
	"net/http"
)

// HandlerFunc 定义所有的请求处理函数 参数是 writer 和 *request
type HandlerFunc func(*Context)

// Engine 实现了 ServeHTTP 接口
type Engine struct {
	*RouterGroup
	router *Router
	groups []*RouterGroup // 存储所有分组
}

// RouterGroup 用于路由分组，支持中间件、分组嵌套
type RouterGroup struct {
	prefix      string        // 路由前缀
	middlewares []HandlerFunc // 支持中间件
	parent      *RouterGroup  // 支持分组的嵌套
	engine      *Engine       // 所有分组共享同一个engine实例
}

// Group 用于在当前分组下创建新的分组
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix:      group.prefix + prefix,
		middlewares: group.middlewares,
		parent:      group,
		engine:      engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// New 是 Engine 构造函数
func New() *Engine {
	engine := &Engine{router: NewRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine} // 自引用
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// AddRoute 用于向路由表中添加路由
func (group *RouterGroup) AddRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// GET 用于添加 HTTP 的 GET 请求
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.AddRoute("GET", pattern, handler)
}

// POST 用于添加 HTTP 的 POST 请求
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.AddRoute("POST", pattern, handler)
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
