package handler

import (
	"context"

	"connectrpc.com/connect"
	greetv1 "github.com/ei-sugimoto/logresdb/api/gen/proto/greet/v1"
	"github.com/ei-sugimoto/logresdb/api/internal/ports/in"
)

type GreetHandler struct {
	uc in.GreetService
}

func NewGreetHandler(uc in.GreetService) *GreetHandler {
	return &GreetHandler{uc: uc}
}

func (g *GreetHandler) Greet(ctx context.Context, req *connect.Request[greetv1.GreetRequest]) (*connect.Response[greetv1.GreetResponse], error) {
	name := req.Msg.Name
	if name == "" {
		name = "World"
	}
	return &connect.Response[greetv1.GreetResponse]{
		Msg: &greetv1.GreetResponse{
			Greeting: g.uc.Greet(name),
		},
	}, nil
}
