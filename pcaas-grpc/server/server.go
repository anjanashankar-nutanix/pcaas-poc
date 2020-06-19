package main

import (
	"context"
	"fmt"
	"pcaas-grpc/pb"
	"net"
	"log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type server struct{}

func (*server) HelloWorld(ctx context.Context,
	req *pb.HelloWorldRequest) (*pb.HelloWorldResponse, error) {
	fmt.Printf("Hello invoked with %v\n", req)
	result := "Hello World"
	res := &pb.HelloWorldResponse {
		Response: result,
	}
	return res, nil
}

func main() {
	fmt.Println("Hello world from Server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	tls := false
	if tls {
		certFile := "ssl/server.crt"
		keyFile := "ssl/server.pem"
		creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
		if sslErr != nil {
			log.Fatalf("Failed loading certificates: %v", sslErr)
			return
		}
		opts = append(opts, grpc.Creds(creds))
	}

	s := grpc.NewServer(opts...)
	pb.RegisterHelloWorldServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
