package main

import (
	"context"
	"github.com/velancio/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

type server struct{}

func (s server) GreetEveryone(everyoneServer greetpb.GreetService_GreetEveryoneServer) error {
	for {
		req, err := everyoneServer.Recv()
		if err == io.EOF{
			return nil
		}
		if err != nil {
			log .Fatalf("Couldn't receive bi di stream: %v", err)
		}
		firstName := req.GetGreeting().GetFirstName()
		result := "Hello " + firstName + "!"
		err = everyoneServer.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})
		if err != nil {
			log .Fatalf("Couldn't send to bi di stream: %v", err)
		}
	}
}

func (s server) LongGreet(greetServer greetpb.GreetService_LongGreetServer) error {
	result := ""
	for {
		req, err:= greetServer.Recv()
		if err == io.EOF{
			return greetServer.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Failed to Recieve from client: %v", err)
		}

		firstName := req.GetGreeting().GetFirstName()
		result = "Hello " + firstName + "!"
	}
}

func (s server) GreetManyTimes(request *greetpb.GreetManyTimesRequest, timesServer greetpb.GreetService_GreetManyTimesServer) error {
	firstName := request.GetGreeting().GetFirstName()
	for i:=0; i<=10; i++ {
		result := "Hello " + firstName + " " + strconv.Itoa(i)
		res := greetpb.GreetManyTimesResponse{
			Result: result,
		}
		err := timesServer.Send(&res)
		if err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (s server) Greet(ctx context.Context, request *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	firstName := request.GetGreeting().GetFirstName()
	lastName := request.GetGreeting().GetLastname()
	if firstName == "<html>"{
		return nil, status.Errorf(codes.InvalidArgument, "Received html tag")
	}
	resp := greetpb.GreetResponse{
		Result: firstName + " " +lastName,
	}
	return &resp, nil
}

func main(){
	lis, err:= net.Listen("tcp", "0.0.0.0:50051")
	if err != nil{
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
}


