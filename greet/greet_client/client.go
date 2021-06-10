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
)

func main(){
	cc, err := grpc.Dial("localhost:50051",grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v",err)
	}
	defer func(cc *grpc.ClientConn) {
		err := cc.Close()
		if err != nil {
			log.Fatalf("Could not close: %v",err)
		}
	}(cc)

	c := greetpb.NewGreetServiceClient(cc)
	doUnary(c)
	doErrorUnary(c)
	doServerStreaming(c)
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
		log.Fatalf("Erro while calling greet rpc: %v", err)
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
		if ok{
			fmt.Println(respErr.Code())
			fmt.Println(respErr.Message())
			if respErr.Code() == codes.InvalidArgument{
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