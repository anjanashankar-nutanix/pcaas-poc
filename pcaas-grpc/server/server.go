package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"log"
	"net"
	"pcaas-grpc/pb"
	"strconv"
	"time"
)

type server struct{}

//Unary RPC implementation
func (*server) HelloWorld(ctx context.Context,
	req *pb.HelloWorldRequest) (*pb.HelloWorldResponse, error) {
	fmt.Printf("Hello invoked with %v\n", req)
	result := "Hello World"
	res := &pb.HelloWorldResponse{
		Response: result,
	}
	return res, nil
}

// Server Streaming RPC implementation
func (*server) HelloWorldServerStream(req *pb.HelloWorldRequest, stream pb.HelloWorldService_HelloWorldServerStreamServer) error {
	fmt.Printf("HelloWorldServerStream function was invoked with %v\n", req)
	message := req.GetHello().GetMessage()
	for i := 0; i < 10; i++ {
		result := message + " number " + strconv.Itoa(i)
		res := &pb.HelloWorldResponse{
			Response: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

//Client Stream RPC implementation
func (*server) HelloWorldClientStream(
	stream pb.HelloWorldService_HelloWorldClientStreamServer) error {
	fmt.Printf("ClientStream function was invoked with a streaming request\n")
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// we have finished reading the client stream
			return stream.SendAndClose(&pb.HelloWorldResponse{
				Response: result,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		message := req.GetHello().GetMessage()
		result += message+ "! "
	}
}

//BiDirectional Streaming
func (*server) HelloWorldBiDirectionalStream(
	stream pb.HelloWorldService_HelloWorldBiDirectionalStreamServer) error {
	fmt.Printf("BiDirectional Stream function was invoked with a streaming" +
		" request\n")

	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}
		message := req.GetHello().GetMessage()
		result += message+ "! "

		sendErr := stream.Send(&pb.HelloWorldResponse{
			Response: result,
		})
		if sendErr != nil {
			log.Fatalf("Error while sending data to client: %v", sendErr)
			return sendErr
		}
	}
}

func main() {
	fmt.Println("Hello world from Server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	tls := true
	if tls {
		certFile := "ssl/service.pem"
		keyFile := "ssl/service.key"
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
