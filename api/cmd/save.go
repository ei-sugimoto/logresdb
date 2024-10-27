package cmd

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"connectrpc.com/grpcreflect"
	"github.com/ei-sugimoto/logresdb/api/gen/proto/greet/v1/greetv1connect"
	"github.com/ei-sugimoto/logresdb/api/internal/adapters/handler"
	"github.com/ei-sugimoto/logresdb/api/internal/adapters/repository"
	"github.com/ei-sugimoto/logresdb/api/internal/usecase"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func Save() {
	db, err := repository.NewDB()
	if err != nil {
		panic(err)
	}
	slog.Info("Connected to database")
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

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, os.Kill)
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	slog.Info("Shutting down")
	defer func() {
		db.Close()
		stop()
		cancel()
	}()

}
