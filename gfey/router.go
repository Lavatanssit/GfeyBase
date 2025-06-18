package gfey

import (
	"log"
	"net/http"
	"strings"
)

// Router: 路由类，使用 Trie 存储 url 到处理句柄的映射。
// 每种 HTTP 请求码对应各自的 TRIE 树
type Router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

// NewRouter: 构造 Router ，返回指向新 Router 实例的指针
func NewRouter() *Router {
	return &Router{
		roots:    make(map[string]*node),
		handlers: map[string]HandlerFunc{},
	}
}

// parsePattern: 将完整的 pattern 解析为 parts 切片，遇见 * 通配符则只解析第一个含 * 的 part
func parsePattern(pattern string) []string {
	// part 切片
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}

	return parts
}

// addRoute: 向 router 中添加路由项
func (r *Router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	_, ok := r.roots[method]

	if !ok {
		// 不存在该路由对应的 method 的情况
		r.roots[method] = &node{}
	}

	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler

	log.Printf("Route %4s - %s", method, pattern)
}

// getRoute: 搜索并返回用户请求 path 的对应 pattern 节点和参数、通配的解析切片。
// 注意 path 是用户 url 路径，pattern 是路由表注册的路由规则
func (r *Router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)

	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil {
		// 解析根节点存的 pattern
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				// e.g. params["userId"] = "admin", part[0:]=":userId", part[1:0]="userId"
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

// Handle: 在路由表中查找上下文请求，解析path参数并放入上下文，并执行对应的handleFunc。
// 查找不到返回404
func (r *Router) Handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)

	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(ctx *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND:%v\n", c.Path)
		})
	}
	c.Next()
}
