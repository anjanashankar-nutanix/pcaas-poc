package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"pcaas-grpc/pb"
	"context"
)

func main() {
	fmt.Printf("Hello I am a client")

	cc, err := grpc.Dial("10.4.206.175:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect : %v", err)
	}

	defer cc.Close()

	c := pb.NewHelloWorldServiceClient(cc)
	req := &pb.HelloWorldRequest{
		Hello : &pb.HelloWorld{
			Message: "Hello world",
		},
	}
	res, err := c.HelloWorld(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling HelloWorld RPC : %v", err)
	}

	fmt.Printf("Response From Server : %v", res.Response)
}

