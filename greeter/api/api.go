package main

import (
	"log"
	"github.com/micro/go-micro"
	"greeter/api/controller/say"
	"greeter/api/controller/order"
	hello "greeter/api/proto/sayhello"
	sayto "greeter/api/proto/saytwo"
)

func main() {
	//定义服务名字
	service := micro.NewService(
		micro.Name("go.micro.api.greeter"),
	)

	// parse command line flags
	service.Init()
	//定义远程调用
	service.Server().Handle(
		service.Server().NewHandler(
			&say.Say{Client: hello.NewSayService("go.micro.srv.greeter", service.Client())},
			//&Order{Client:  sayto.NewOrderService("go.micro.srv.greeter", service.Client())},
		),
	)
	//定义使用路由 这里应该有简单方法 暂时未找到 先这样吧
	service.Server().Handle(
		service.Server().NewHandler(
			//&Say{Client: hello.NewSayService("go.micro.srv.greeter", service.Client())},
			&order.Order{Client: sayto.NewOrderService("go.micro.srvtwo.greeter", service.Client())},
		),
	)
	//启动服务
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
