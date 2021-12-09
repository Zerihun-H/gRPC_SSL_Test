package main

import (
	"context"
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

	creds, err := credentials.NewServerTLSFromFile("/etc/letsencrypt/live/adinterface.adsrecognition.com-0002/fullchain.pem;", "/etc/letsencrypt/live/adinterface.adsrecognition.com-0002/privkey.pem")
	if err != nil {
		log.Fatal("msg", "failed to setup TLS with local files", "error", err)

	}

	opts := []grpc.ServerOption{
		// Enable TLS for all incoming connections.
		grpc.Creds(creds),
	}

	lis, err := net.Listen("tcp", ":8383")
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
