package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	example "github.com/hatlonely/hellogrpc/go/api/gen/go/api"
)

type EchoService struct{}

func (s *EchoService) Echo(ctx context.Context, req *example.EchoReq) (*example.EchoRes, error) {
	return &example.EchoRes{Value: req.Value}, nil
}

type CalService struct{}

func (s *CalService) Cal(ctx context.Context, req *example.CalReq) (*example.CalRes, error) {
	var result int64
	switch req.Info.Op {
	case "+":
		result = req.Info.A + req.Info.B
	case "-":
		result = req.Info.A - req.Info.B
	default:
		return nil, status.Errorf(codes.InvalidArgument, "op should in ['+', '-']")
	}

	return &example.CalRes{
		Result: result,
		Uid:    req.Uid,
	}, nil
}

func main() {
	if len(os.Args) > 1 {
		server := grpc.NewServer()
		example.RegisterEchoServiceServer(server, &EchoService{})
		example.RegisterCalServiceServer(server, &CalService{})
		address, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", 6060))
		if err != nil {
			panic(err)
		}
		if err := server.Serve(address); err != nil {
			panic(err)
		}
	} else {
		mux := runtime.NewServeMux()
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		if err := example.RegisterEchoServiceHandlerServer(ctx, mux, &EchoService{}); err != nil {
			panic(err)
		}
		if err := example.RegisterCalServiceHandlerServer(ctx, mux, &CalService{}); err != nil {
			panic(err)
		}
		if err := http.ListenAndServe(":80", mux); err != nil {
			panic(err)
		}
	}
}
