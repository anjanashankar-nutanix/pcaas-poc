package main

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"pcaas-grpc/pb"
	"context"
)

func main() {
	fmt.Printf("Hello I am a client")

	/*cc, err := grpc.Dial("10.4.206.175:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect : %v", err)
	}

	defer cc.Close()*/

	tls := true
	opts := grpc.WithInsecure()
	if tls {
		certFile := "ssl/ca.cert" // Certificate Authority Trust certificate
		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
		if sslErr != nil {
			log.Fatalf("Error while loading CA trust certificate: %v", sslErr)
			return
		}
		opts = grpc.WithTransportCredentials(creds)
	}

	cc, err := grpc.Dial("10.138.102.133:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
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

