package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"log"
	"pcaas-grpc/pb"
	"time"
)

func main() {
	fmt.Printf("Hello I am a client\n")
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
	// Calling the Unary RPC
	// doUnary(c)
	//doServerStream(c)
	//doClientStream(c)
	doBiDirectionalStream(c)
}

func doBiDirectionalStream(c pb.HelloWorldServiceClient) {
	fmt.Println("Starting to do a BiDi Streaming RPC...")

	// we create a stream by invoking the client
	stream, err := c.HelloWorldBiDirectionalStream(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	requests := []*pb.HelloWorldRequest{
		&pb.HelloWorldRequest{
			Hello: &pb.HelloWorld{
				Message: "Hello world 1",
			},
		},
		&pb.HelloWorldRequest{
			Hello: &pb.HelloWorld{
				Message: "Hello world 2",
			},
		},
		&pb.HelloWorldRequest{
			Hello: &pb.HelloWorld{
				Message: "Hello world 3",
			},
		},
	}

	waitc := make(chan struct{})
	// we send a bunch of messages to the client (go routine)
	go func() {
		// function to send a bunch of messages
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	// we receive a bunch of messages from the client (go routine)
	go func() {
		// function to receive a bunch of messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetResponse())
		}
		close(waitc)
	}()

	// block until everything is done
	<-waitc
}

func doClientStream(c pb.HelloWorldServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")
	requests := []*pb.HelloWorldRequest{
		&pb.HelloWorldRequest{
			Hello: &pb.HelloWorld{
				Message: "Hello world 1",
			},
		},
		&pb.HelloWorldRequest{
			Hello: &pb.HelloWorld{
				Message: "Hello world 2",
			},
		},
		&pb.HelloWorldRequest{
			Hello: &pb.HelloWorld{
				Message: "Hello world 3",
			},
		},
	}
	stream, err := c.HelloWorldClientStream(context.Background())
	if err != nil {
		log.Fatalf("error while calling ClientStream: %v", err)
	}
	// we iterate over our slice and send each message individually
	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from ClientStream: %v", err)
	}
	fmt.Printf("ClientStream Response: %v\n", res)
}

func doServerStream(c pb.HelloWorldServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")
	req := &pb.HelloWorldRequest{
		Hello: &pb.HelloWorld{
			Message: "Hello world",
		},
	}
	resStream, err := c.HelloWorldServerStream(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Server Stream RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("Response from Server Stream RPC: %v", msg.GetResponse())
	}
}

func doUnary(c pb.HelloWorldServiceClient) {
	req := &pb.HelloWorldRequest{
		Hello: &pb.HelloWorld{
			Message: "Hello world",
		},
	}
	res, err := c.HelloWorld(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling HelloWorld RPC : %v", err)
	}
	fmt.Printf("Response From Server : %v\n", res.Response)
}
