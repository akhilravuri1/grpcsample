package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "github.com/akhilravuri1/grpcsample/generated"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewCalculatorServiceClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Sum(ctx, &pb.SumRequest{FirstNumber: name, SecondNumber: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	for i, res := range r.GetSumResult() {
		log.Printf("%d : %s", i, res)
	}
}
