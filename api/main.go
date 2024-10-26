package main

import (
	"net/http"

	"connectrpc.com/grpcreflect"
	"github.com/ei-sugimoto/logresdb/api/gen/proto/greet/v1/greetv1connect"
	"github.com/ei-sugimoto/logresdb/api/internal/adapters/handler"
	"github.com/ei-sugimoto/logresdb/api/internal/usecase"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	greetService := usecase.NewGreetService()
	greetHandler := handler.NewGreetHandler(greetService)

	mux := http.NewServeMux()
	path, handler := greetv1connect.NewGreetServiceHandler(greetHandler)
	reflector := grpcreflect.NewStaticReflector(
		"greet.v1.GreetService",
		// protoc-gen-connect-go generates package-level constants
		// for these fully-qualified protobuf service names, so you'd more likely
		// reference userv1.UserServiceName and groupv1.GroupServiceName.
	)
	mux.Handle(path, handler)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle((grpcreflect.NewHandlerV1Alpha(reflector)))
	http.ListenAndServe(":8000", h2c.NewHandler(mux, &http2.Server{}))
}
