package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/Lukski175/grpc101/time"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

var stack []string

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Reply: "Hello " + in.GetName()}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) ReceiveMessages(ctx context.Context, in *pb.MessageRequest) (*pb.MessageReply, error) {
	stack = append(stack, in.Message)
	log.Printf("Received message: %v", in.GetMessage())
	return &pb.MessageReply{Messages: nil}, nil
}

func (s *server) SendMessages(ctx context.Context, in *pb.MessageAmount) (*pb.MessageReply, error) {
	var temp []string
	for i := len(stack); i > 0 && i > len(stack)-5; i-- {
		temp = append(temp, stack[i])
	}
	return &pb.MessageReply{Messages: temp}, nil
}
