package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	example "github.com/hatlonely/hellogrpc/go/api/gen/go/api"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:6060", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("dial failed. err: [%v]\n", err)
		return
	}
	defer conn.Close()

	client := example.NewYourServiceClient(conn)

	res, err := client.Echo(context.Background(), &example.StringMessage{Value: "hello world"})
	fmt.Println(res, err)
}
