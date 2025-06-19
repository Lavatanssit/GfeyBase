package main

import (
	"fmt"
	"gfey"
	"gfey/middlewares"
	"net/http"
	"text/template"
	"time"
)

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d:%02d:%02d", year, month, day)
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	// 创建 Engine 实例 root
	root := gfey.New()
	// 使用中间件 Logger
	root.Use(middlewares.Logger())
	root.Use(middlewares.Recovery())
	// 注册模板函数
	root.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	// 加载本地模板文件
	root.LoadHTMLGlob("templates/*")
	// 把本地static目录下的静态资源映射到/assets路径
	root.Static("/assets", "./static")

	// user_admin := &User{
	// 	Username: "admin",
	// 	Password: "admin123456",
	// }
	// user_guest := &User{
	// 	Username: "guest",
	// 	Password: "guest123456",
	// }

	// 注册路由处理函数
	root.GET("/", func(c *gfey.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})

	root.GET("/panic", func(c *gfey.Context) {
		names := []string{"geektutu"}
		c.String(http.StatusOK, names[100])
	})

	root.Run(":10000")
}
