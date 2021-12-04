package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"gRPCTest/credit"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type server struct {
	credit.UnimplementedCreditServiceServer
}

func main() {
	log.Println("Server running ...")

	cert, err := tls.LoadX509KeyPair("ca.cert", "ca.key")
	if err != nil {
		log.Fatalln(err)
	}

	opts := []grpc.ServerOption{
		// Enable TLS for all incoming connections.
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
	}

	lis, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalln(err)
	}

	srv := grpc.NewServer(opts...)
	credit.RegisterCreditServiceServer(srv, &server{})

	log.Fatalln(srv.Serve(lis))
}

func (s *server) Credit(ctx context.Context, request *credit.CreditRequest) (*credit.CreditResponse, error) {
	log.Println(fmt.Sprintf("Request: %g", request.GetAmount()))

	return &credit.CreditResponse{Confirmation: fmt.Sprintf("Credited %g", request.GetAmount())}, nil
}
