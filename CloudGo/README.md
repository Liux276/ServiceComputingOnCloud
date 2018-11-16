## 服务计算 - 3 | golang web编程学习
### CloudGo
---
#### 框架选择

> Iris: 目前发展最快的Go Web框架。提供完整的MVC功能并且面向未来。

ris自称是Go语言中所有Web框架最快的，它的特点如下：
1. 聚焦高性能 
2. 健壮的静态路由支持和通配符子域名支持。 
3. 视图系统支持超过5以上模板 
4. 支持定制事件的高可扩展性Websocket API 
5. 带有GC, 内存 & redis 提供支持的会话 
6. 方便的中间件和插件 
7. 完整 REST API 
8. 能定制 HTTP 错误 
9. Typescript编译器 + 基于浏览器的编辑器 
10. 内容 negotiation & streaming 
11. 传送层安全性 
12. 源码改变后自动加载 
13. OAuth, OAuth2 支持27+ API providers 
14. JSON Web Tokens

Iris是目前最优秀的框架之一，提供完整的WEB组件，且拥有较高的运行速度，所以我选择Iris进行此次开发。

---
#### 服务器代码

```go
package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

func main() {
	app := iris.New()
	app.Use(recover.New())
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
	app.Run(iris.Addr(":9090")) //9090 监听端口
}

```

---
#### 服务器测试
1. 运行服务器：
![运行服务器](https://segmentfault.com/img/bVbjD8c?w=558&h=104)
2. 浏览器访问：
![浏览器访问](https://segmentfault.com/img/bVbjD8e?w=467&h=209)
3. 使用curl对可以正常访问服务器：
![Curl访问结果](https://segmentfault.com/img/bVbjD8s?w=582&h=334)

---
#### AB测试

* 测试环境：centos7
* 测试请求数：10000
* 测试命令：`ab -n 10000 -c 100 http://localhost:9090/welcome`
* 测试结果：
![AB测试](https://segmentfault.com/img/bVbjD77?w=659&h=704)
* 测试结果重要参数：

|字段|含义|
| :---: | :---: |
|Server Hostname|服务器主机名|
|Server Port|服务器端口|
|Document Path|文件路径|
|Document Length|文件大小|
|Concurrency Level|并发等级|
|Requst per second|平均每秒的请求个数。服务器并发处理能力的量化描述，单位是reqs/s，指的是在某个并发用户数下单位时间内处理的请求数。某个并发用户数下单位时间内能处理的最大请求数，称之为最大吞吐率。<br>吞吐率是基于并发用户数的。这句话代表了两个含义：<br>a、吞吐率和并发用户数相关。<br>b、不同的并发用户数下，吞吐率一般是不同的。<br>计算公式：总请求数/处理完成这些请求数所花费的时间，即：Request per second=Complete requests/Time taken for tests.这个数值表示当前机器的整体性能，值越大越好。|
|Time per request|用户平均的等待时间。<br>计算公式：总请求数/处理完成这些请求数所花费的时间，即：Request per second=Complete requests/Time taken for tests.这个数值表示当前机器的整体性能，值越大越好。|
|Time per request:across all concurrent requests|计算公式：处理完成所有请求数所花费的时间/总请求数，即：Time taken for/testsComplete requests.<br>可以看到，它是吞吐率的倒数。同时，它也等于用户平均请求等待时间/并发用户数，即Time per request/Concurrency Level。|
|Connection Times|表内描述了所有的过程中所消耗的最小、中位、最长时间。|
|Percentage of the requests served within a certain time|每个百分段的请求完成所需的时间|

---
### net/http源码阅读

[NET/HTTP源码阅读](https://github.com/Liux276/ServiceComputingOnCloud/blob/master/CloudGo/httpSourceCodeReading.md)

---
