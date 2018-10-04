package order

import (
	sayto "greeter/api/proto/saytwo"

	"log"
	"github.com/micro/go-micro/errors"
	"strings"
	"encoding/json"
	"context"
	api "github.com/micro/go-api/proto"
)

type Order struct {
	Client sayto.OrderService
}

func (s *Order) Yes(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Received Order.Yes API request")

	name, ok := req.Get["name"]
	if !ok || len(name.Values) == 0 {
		return errors.BadRequest("go.micro.api.greeter", "Name cannot be blank")
	}

	response, err := s.Client.Yes(ctx, &sayto.Request{
		Name: strings.Join(name.Values, " "),
	})
	if err != nil {
		return err
	}

	rsp.StatusCode = 200
	b, _ := json.Marshal(map[string]string{
		"message": response.Msg,
	})
	rsp.Body = string(b)

	return nil
}


func (s *Order) No(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Received Order.No API request")

	name, ok := req.Get["name"]
	if !ok || len(name.Values) == 0 {
		return errors.BadRequest("go.micro.api.greeter", "Name cannot be blank")
	}

	response, err := s.Client.No(ctx, &sayto.Request{
		Name: strings.Join(name.Values, " "),
	})
	if err != nil {
		return err
	}

	rsp.StatusCode = 200
	b, _ := json.Marshal(map[string]string{
		"message": response.Msg,
	})
	rsp.Body = string(b)

	return nil
}