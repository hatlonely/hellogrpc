package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	example "github.com/hatlonely/hellogrpc/go/api/gen/go/api"
)

type Service struct{}

func (s *Service) Echo(ctx context.Context, req *example.StringMessage) (*example.StringMessage, error) {
	return &example.StringMessage{Value: req.Value}, nil
}

func main() {
	go func() {
		server := grpc.NewServer()
		example.RegisterYourServiceServer(server, &Service{})
		address, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", 6060))
		if err != nil {
			panic(err)
		}
		if err := server.Serve(address); err != nil {
			panic(err)
		}
	}()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if err := example.RegisterYourServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("0.0.0.0:%v", 6060), opts); err != nil {
		panic(err)
	}
	if err := http.ListenAndServe(":80", mux); err != nil {
		panic(err)
	}
}
