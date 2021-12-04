package main

import (
	"context"
	"crypto/tls"
	pb "gRPCTest/helloworld"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

func main() {
	log.Println("Server running ...")

	cert, err := tls.LoadX509KeyPair("service.pem", "service.key")
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
	pb.RegisterGreeterServer(srv, &server{})

	log.Fatalln(srv.Serve(lis))
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}
