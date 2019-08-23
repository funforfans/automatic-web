package main

import (
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro/web"
	"net/http"
	"tyweiqu/handler"
	"github.com/micro/cli"
)

func main() {

	// 创建新服务
	service := web.NewService(
		// 后面两个web，第一个是指是web类型的服务，第二个是服务自身的名字
		web.Name("tuyoo.micro.web.tools"),
		web.Version("latest"),
		web.Address(":63000"),
	)

	// 初始化服务
	if err := service.Init(
		web.Action(
			func(c *cli.Context) {
				// 初始化handler
				handler.Init()
			}),
		); err != nil {
		log.Fatal(err)
	}
	// 这是返回一个静态网页，如果直接建立连接可以忽略
	service.Handle("/", http.FileServer(http.Dir("html")))
	// 注册登录接口
	service.HandleFunc("/trigger", handler.Trigger)
	service.HandleFunc("/upload", handler.UploadHandler)
	service.HandleFunc("/view", handler.ViewHandler)
	service.HandleFunc("/webhook", handler.WebhookHandler)
	// 运行服务
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
