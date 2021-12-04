package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	pb "gRPCTest/helloworld"
	"io/ioutil"
	"log"
	"time"

	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
)

const (
	Addr       = ":50052"
	AkbarCA    = "cert/adinterface.adsrecognition.com.crt"
	ClientCert = "cert/client.crt"
	ClientKey  = "cert/client.key"
	ServerName = "server"
)

func main() {
	log.Println("Client running ...")

	// Load certs from the d
	cert, err := tls.LoadX509KeyPair(ClientCert, ClientKey)
	if err != nil {
		fmt.Errorf("Could not load client key pair : %v", err)
	}

	// Create certpool from the CA
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(AkbarCA)
	if err != nil {
		fmt.Errorf("Could not read Cert CA : %v", err)
	}

	// Append the certs from the CA
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		fmt.Errorf("Failed to append CA cert : %v", err)
	}

	// Create transport creds based on TLS.
	creds := credentials.NewTLS(&tls.Config{
		ServerName:   ServerName,
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
	})

	opts := []grpc.DialOption{
		// grpc.WithInsecure(),
		grpc.WithTransportCredentials(creds),
		grpc.WithBlock(),
	}

	conn, err := grpc.Dial("operatorsvps2.adsrecognition.com:50051", opts...)
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
