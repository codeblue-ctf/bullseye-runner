package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc/credentials"

	pb "gitlab.com/CBCTF/bullseye-runner/proto"
	"google.golang.org/grpc"
)

var (
	certFile = flag.String("cert", "", "TLS certificate")
	keyFile  = flag.String("key", "", "TLS private key")
	port     = flag.Int("port", 10080, "port to listen")
)

type runnerServer struct{}

func (s *runnerServer) Run(ctx context.Context, req *pb.RunnerRequest) (*pb.RunnerResponse, error) {
	log.Printf("received")

	time.Sleep(3 * time.Second)

	log.Printf("uuid: %s", req.Uuid)

	res := pb.RunnerResponse{
		Uuid:      req.Uuid,
		Succeeded: true,
		Stdout:    "stdout",
		Stderr:    "stderr",
	}
	return &res, nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	if *certFile != "" && *keyFile != "" {
		creds, err := credentials.NewClientTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	server := grpc.NewServer(opts...)
	pb.RegisterRunnerServer(server, &runnerServer{})
	server.Serve(lis)
}
