package main

import (
	"context"
	"gRPCTest/credit"
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

	conn, err := grpc.Dial(":50051", opts...)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	client := credit.NewCreditServiceClient(conn)

	request := &credit.CreditRequest{Amount: 1990.01}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.Credit(ctx, request)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Response:", response.GetConfirmation())
}
