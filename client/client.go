package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/Lukski175/grpc101/time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", "something", "Name to greet")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

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
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetReply())

	go Chat()

	for {
		scan := bufio.NewScanner(os.Stdin)
		scan.Scan()
		input := scan.Text()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_, err := c.ReceiveMessages(ctx, &pb.MessageRequest{Message: input})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
	}
}

func Chat(c *pb.GreeterClient) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := c.SendMessages(ctx, &pb.MessageAmount{amount: 10})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		_ = r.GetReply()
	}
}
