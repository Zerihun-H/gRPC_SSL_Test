package main

import (
	"context"
	pb "gRPCTest/helloworld"
	"log"
	"time"

	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
)

const (
	AkbarCA    = "/var/www/SSL/adinterface.adsrecognition.com.crt"
	ClientCert = "/var/www/SSL/ca.cert"
	ClientKey  = "/var/www/SSL/client.key"
	ServerName = "server"
)

func main() {
	log.Println("Client running ...")
	creds, err := credentials.NewClientTLSFromFile("/etc/letsencrypt/live/adinterface.adsrecognition.com-0002/fullchain.pem;", "localhost")
	if err != nil {
		log.Fatalln(err)
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithBlock(),
	}

	conn, err := grpc.Dial("adinterface.adsrecognition.com:50051", opts...)
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
