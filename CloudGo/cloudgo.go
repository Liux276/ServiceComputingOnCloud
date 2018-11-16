package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

func main() {
	//创建新iris对象
	app := iris.New()
	//恢复方法
	app.Use(recover.New())
	//logger输出器
	app.Use(logger.New())
	//输出html
	// 请求方式: GET
	// 访问地址: http://localhost:9090/welcome
	app.Handle("GET", "/welcome", func(ctx iris.Context) {
		ctx.HTML("<h1>Welcome</h1>")
	})
	//输出字符串
	// 类似于 app.Handle("GET", "/ping", [...])
	// 请求方式: GET
	// 请求地址: http://localhost:9090/ping
	app.Get("/ping", func(ctx iris.Context) {
		ctx.WriteString("pong")
	})
	//输出json
	// 请求方式: GET
	// 请求地址: http://localhost:9090/hello
	app.Get("/hello", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"message": "Hello World!"})
	})
	//监听9090端口
	app.Run(iris.Addr(":9090")) //9090 监听端口
}
