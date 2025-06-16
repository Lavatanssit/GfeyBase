package gfey

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

// Context: *http.Request 和 http.ResponseWriter 的封装
type Context struct {
	// 核心字段
	Writer http.ResponseWriter
	Req    *http.Request
	// 请求信息字段
	Path   string
	Method string
	Params map[string]string
	// 响应信息字段
	StatusCode int
}

// NewContext: 构造Context，返回指向新Context实例的指针
func NewContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer:     w,
		Req:        req,
		Path:       req.URL.Path,
		Method:     req.Method,
		StatusCode: 0,
	}
}

// PostForm: 根据key，提取HTTP-POST请求体中的Form参数
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query：根据key，提取GET-url中?后的Query参数
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// Status：设置HTTP响应的状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader: 设置HTTP响应的头部字段（比如内容类型JSON、HTML或是其他自定义字段）
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// String: 构造String响应的方法，values可变接口参数，用于填充到format字符串内的占位符
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// Json: 构造Json响应的方法，使用encoder绑定writer，对接口jsonObj类型Encode时自动写入writer
func (c *Context) Json(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

// Data: 构造Data响应的方法，向writer中写入data数据
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// HTML: 构造HTML响应的方法，向writer中写入string类型的html页面
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

func (c *Context) Param(key string) string {
	value := c.Params[key]
	return value
}
