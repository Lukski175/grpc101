package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	pb "github.com/Lukski175/grpc101/time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func portMethod() (port *int) {
	return flag.Int("port", porten, "The server port")
}

var messages []*pb.ClientMessage

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.GreeterServer
}

var porten int

func tupleMethod() (addr *string, name *string) {
	fmt.Printf(strconv.Itoa(porten))
	return flag.String("addr", "localhost:"+strconv.Itoa(porten), "the address to connect to"),
		flag.String("name", "something", "Name to greet")
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	porten++
	go ServerSetup()
	return &pb.HelloReply{Reply: "Hello " + in.GetName(), Port: int32(porten)}, nil
}

var gs pb.GreeterServer
var gc pb.GreeterClient

func main() {
	porten = 50051
	//Setup this server for client input
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *portMethod()))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	porten++
}

func (s *server) ReceiveMessages(ctx context.Context, in *pb.MessageRequest) (*pb.MessageReply, error) {
	messages = append(messages, in.GetMessage())
	log.Printf("Received message: %v", in.GetMessage())
	SendMessages()
	return &pb.MessageReply{Messages: nil}, nil
}

func ServerSetup() {
	//Setup server connection to client, so server can send messages
	address, _ := tupleMethod()
	conn, err := grpc.Dial(*address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	gc = pb.NewGreeterClient(conn)
}

func SendMessages() {
	//log.Print("Sending:", elem)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := gc.Chat(ctx, &pb.MessageReply{Messages: messages})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
}
