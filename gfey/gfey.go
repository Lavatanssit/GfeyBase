package gfey

import (
	"log"
	"net/http"
	"path"
	"strings"
	"text/template"
)

// HandlerFunc 定义所有的请求处理函数 参数是 *Context
type HandlerFunc func(*Context)

// Engine 实现了 ServeHTTP 接口
type Engine struct {
	*RouterGroup
	router        *Router
	groups        []*RouterGroup     // 存储所有分组
	htmlTemplates *template.Template // html渲染模板
	funcMap       template.FuncMap   // html渲染函数
}

// New 是 Engine 构造函数
func New() *Engine {
	engine := &Engine{router: NewRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine} // 自引用
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// RouterGroup 用于路由分组，支持多级嵌套和中间件链
type RouterGroup struct {
	prefix      string        // 路由前缀
	middlewares []HandlerFunc // 支持中间件
	parent      *RouterGroup  // 支持分组的嵌套
	engine      *Engine       // 所有分组共享同一个engine实例
}

// Group 用于在当前分组下创建新的分组，返回新分组指针，便于链式调用
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

// USE 用于向当前分组添加中间件，便于链式调用
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// Run 用于运行 http 服务器，调用标准库 ListenAndServe 方法
func (engine *Engine) Run(address string) (err error) {
	return http.ListenAndServe(address, engine)
}

// Engine 实现标准库的 ServeHTTP 接口，自动收集匹配分组的中间件并执行
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 给对应分组增加中间件
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := NewContext(w, req)
	c.handlers = middlewares
	c.engine = engine
	engine.router.Handle(c)
}

// createStaticHandler 创建静态文件处理函数。
// relativePath: 路由分组下的静态资源访问前缀（如 "/" 或 "/assets"）。
// fs: 实现 http.FileSystem 接口的文件系统（如 http.Dir）。
// 返回值：用于处理静态资源请求的 HandlerFunc。
// 处理流程：
//  1. 计算绝对路由前缀 absolutePath（如 "/assets"）。
//  2. 用 http.StripPrefix 去除 URL 路径前缀，保证本地目录和 URL 对应。
//  3. 检查文件是否存在，不存在返回 404，存在则交给标准库 fileServer 处理。
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(ctx *Context) {
		file := ctx.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			ctx.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(ctx.Writer, ctx.Req)
	}
}

// Static 注册静态文件路由，将 URL 路径映射到本地文件系统。
// relativePath: 路由分组下的静态资源访问前缀（如 "/" 或 "/assets"）。
// root: 本地静态文件根目录（如 "./static"）。
// 例如：group.Static("/assets", "./static")，则访问 /assets/css/style.css 会映射到 ./static/css/style.css。
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")

	group.GET(urlPattern, handler)
}

// SetFuncMap 设置模板渲染时可用的自定义函数。
// 注意：应在 LoadHTMLGlob 之前调用。
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

// LoadHTMLGlob 加载模板文件，支持通配符（如 "templates/*"）。
// pattern: 模板文件路径模式。
// 加载后可通过 Context.HTML 方法渲染模板页面。
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}
