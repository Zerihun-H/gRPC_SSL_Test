package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	pb "gRPCTest/helloworld"
	"io/ioutil"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Cert file details
const (
	AkbarCA    = "cert/adinterface.adsrecognition.com.crt"
	ServerCert = "cert/server.crt"
	ServerKey  = "cert/server.key"
	ServerName = "server"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

func main() {
	log.Println("Server running ...")

	// Load certs from the disk.
	cert, err := tls.LoadX509KeyPair(ServerCert, ServerKey)
	if err != nil {
		fmt.Errorf("could not server key pairs: %s", err)
	}

	// Create certpool from the CA
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(AkbarCA)
	if err != nil {
		fmt.Printf("could not read CA cert: %s", err)
	}

	// Append the certs from the CA
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		fmt.Printf("Failed to append the CA certs: %s", err)
	}

	// Create the TLS config for gRPC server.
	creds := credentials.NewTLS(
		&tls.Config{
			ClientAuth:   tls.RequireAnyClientCert,
			Certificates: []tls.Certificate{cert},
			ClientCAs:    certPool,
		})
	opts := []grpc.ServerOption{
		// Enable TLS for all incoming connections.
		grpc.Creds(creds),
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
