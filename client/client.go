package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"

	pb "github.com/Lukski175/grpc101/time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var porten int

func tupleMethod() (addr *string, name *string) {
	return flag.String("addr", "localhost:"+strconv.Itoa(porten), "the address to connect to"),
		flag.String("name", "something", "Name to greet")
}

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.GreeterServer
}

func portMethod() (port *int) {
	return flag.Int("port", porten, "The server port")
}

func main() {
	porten = 50051
	flag.Parse()
	// Set up a connection to the server.
	address, name := tupleMethod()
	conn, err := grpc.Dial(*address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	//Sets up client name
	fmt.Print("Enter client name:")
	scan := bufio.NewScanner(os.Stdin)
	scan.Scan()
	input := scan.Text()
	name = &input

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Greeting: %s", r.GetReply())
	porten = int(r.GetPort())

	go AwaitServer()

	//Send text to server loop
	for {
		scan := bufio.NewScanner(os.Stdin)
		scan.Scan()
		input := scan.Text()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_, err := c.ReceiveMessages(ctx, &pb.MessageRequest{Message: &pb.ClientMessage{Name: *name, Message: input}})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
	}
}

func AwaitServer() {
	//Sets up server to receive chat messages
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
}

func (s *server) Chat(ctx context.Context, in *pb.MessageReply) (*pb.HelloRequest, error) {

	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()

	for _, v := range in.GetMessages() {
		fmt.Println(v.Name + ": " + v.Message)
	}

	return &pb.HelloRequest{Name: ""}, nil
}
