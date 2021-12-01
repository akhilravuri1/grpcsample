package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/akhilravuri1/grpcsample/generated"
	"github.com/couchbase/gocb/v2"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedCalculatorServiceServer
}

func (s *server) Sum(ctx context.Context, in *pb.SumRequest) (*pb.SumResponse, error) {
	Cluster, err := gocb.Connect("localhost",
		gocb.ClusterOptions{
			Username: "sample",
			Password: "sample",
		})
	if err != nil {
		fmt.Println("ERRROR CONNECTING TO CLUSTER:", err)
	}
	// Open Bucket
	Bucket := Cluster.Bucket("sample-bucket")
	err = Bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		fmt.Println("Error while waiting : ", err)
	}
	res, err := Cluster.Query("select * from `sample-bucket`.`sample-scope`.`sample-sellers`", nil)
	if err != nil {
		fmt.Println("Error while reteriving content Using Query: ", err)
	}

	var SellersDb []*pb.Seller
	for res.Next() {
		var temp map[string]*pb.Seller
		err := res.Row(&temp)
		if err != nil {
			fmt.Println("This is Row Error in GetAllSams: ", err)
		}
		SellersDb = append(SellersDb, temp["sample-sellers"])
	}
	fmt.Println(SellersDb)
	//a := []string{"akhil, nikhil"}
	return &pb.SumResponse{SumResult: SellersDb}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCalculatorServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
