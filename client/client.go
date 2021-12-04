package main

import (
	"context"
	pb "gRPCTest/helloworld"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	log.Println("Client running ...")

	creds, err := credentials.NewClientTLSFromFile("../ca.cert", "operatorsvps2.adsrecognition.com")
	if err != nil {
		log.Fatalln(err)
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithBlock(),
	}

	conn, err := grpc.Dial(":8888", opts...)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.SayHello(ctx, &pb.HelloRequest{Name: "Zerihun"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())

}
