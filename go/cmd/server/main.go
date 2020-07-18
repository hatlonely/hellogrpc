package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/peer"
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
		server := grpc.NewServer(
			grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res interface{}, err error) {
				ts := time.Now()
				p, ok := peer.FromContext(ctx)
				clientIP := ""
				if ok && p != nil {
					clientIP = p.Addr.String()
				}
				res, err = handler(ctx, req)
				fmt.Println(map[string]interface{}{
					"client":    clientIP,
					"url":       info.FullMethod,
					"req":       req,
					"res":       res,
					"err":       err,
					"resTimeNs": time.Now().Sub(ts).Nanoseconds(),
				})
				return res, err
			}),
			grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
				MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
				PermitWithoutStream: true,            // Allow pings even when there are no active streams
			}),
			grpc.KeepaliveParams(keepalive.ServerParameters{
				MaxConnectionIdle:     15 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
				MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
				MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
				Time:                  5 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
				Timeout:               1 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
			}),
		)
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
