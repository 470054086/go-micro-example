package main

import (
	"log"
	"time"
	"github.com/micro/go-micro"
	"context"
	saywo "greeter/srvtwo/proto/saywo"
)

type Order struct{}

func (s *Order) Yes(ctx context.Context, req *saywo.Request, rsp *saywo.Response) error {
	log.Print("Received Say.Hello request")
	rsp.Msg = "Hello,我是第一个 " + req.Name
	return nil
}

func (s *Order) No(ctx context.Context,req *saywo.Request,res *saywo.Response) error{
	log.Print("Received Say.World request")
	res.Msg = "World" + req.Name
	return nil
}



func main() {

	service := micro.NewService(
		micro.Name("go.micro.srvtwo.greeter"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
	)




	// optionally setup command line usage
	service.Init()
	saywo.RegisterOrderHandler(service.Server(), new(Order))
	// Register Handlers

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}


}
