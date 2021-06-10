package main

import (
	"context"
	"fmt"
	"github.com/velancio/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"time"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer func(cc *grpc.ClientConn) {
		err := cc.Close()
		if err != nil {
			log.Fatalf("Could not close: %v", err)
		}
	}(cc)

	c := greetpb.NewGreetServiceClient(cc)
	/*log.Println("============= GRPC Unary ==============")
	doUnary(c)
	log.Println("============= GRPC Unary with Error ==============")
	doErrorUnary(c)
	log.Println("============= GRPC Server Stream ==============")
	doServerStreaming(c)
	log.Println("============= GRPC Client Stream ==============")
	doClientStreaming(c)*/
	log.Println("============= GRPC BiDi Stream ==============")
	doBiDiStreaming(c)
	log.Println("============= Done ==============")
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while calling bidi streaming greet rpc: %v", err)
	}

	requests := []*greetpb.GreetEveryoneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "John",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Avi",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Shaun",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Rick",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Mike",
			},
		},
	}
	waitc := make(chan struct{})
	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Erro while reading stream: %v", err)
			}
			log.Printf("Response from greeting: %v", msg.Result)
		}
		close(waitc)
	}()

	go func() {
		for _, req := range requests {
			fmt.Printf("Sending the req: %v \n", req)
			err := stream.Send(req)
			if err != nil {
				log.Fatalf("Couldn't send client request stream: %v", err)
			}
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()


	<- waitc
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	requests := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "John",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Avi",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Shaun",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Rick",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Mike",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling client streaming greet rpc: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending the req: %v \n", req)
		err := stream.Send(req)
		if err != nil {
			log.Fatalf("Couldn't send client request stream: %v", err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Couldn't close client request stream: %v", err)
	}

	fmt.Printf("Long greet response: %v \n", res)
}

func doUnary(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "velancio",
			Lastname:  "Fernandes",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling greet rpc: %v", err)
	}
	log.Printf("Response from greeting: %v", res.Result)
}

func doErrorUnary(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "<html>",
			Lastname:  "Fernandes",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		//log.Fatalf("Error while calling greet rpc: %v", err)
		respErr, ok := status.FromError(err)
		if ok {
			fmt.Println(respErr.Code())
			fmt.Println(respErr.Message())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("Html tags are not allowed")
			}
		} else {
			log.Fatalf("Big error:: %v", err)
		}
		return
	}
	log.Printf("Response from greeting: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "velancio",
			Lastname:  "Fernandes",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Erro while calling greet server stream rpc: %v", err)
	}

	for {
		msg, err := resStream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Erro while reading stream: %v", err)
		}

		log.Printf("Response from greeting: %v", msg.Result)
	}
}
