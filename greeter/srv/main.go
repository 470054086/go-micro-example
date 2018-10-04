package main

import (
	"log"
	"time"
	hello "greeter/srv/proto/worldhello"
	"github.com/micro/go-micro"
	"context"
)

type Say struct{}

func (s *Say) Hello(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
	log.Print("Received Say.Hello request")
	rsp.Msg = "Hello,我是第一个 " + req.Name
	return nil
}

func (s *Say) World(ctx context.Context,req *hello.Request,res *hello.Response) error{
	log.Print("Received Say.World request")
	res.Msg = "World" + req.Name
	return nil
}



func main() {
	service := micro.NewService(
		micro.Name("go.micro.srv.greeter"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
	)
	// optionally setup command line usage
	service.Init()
	// Register Handlers
	hello.RegisterSayHandler(service.Server(), new(Say))

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
