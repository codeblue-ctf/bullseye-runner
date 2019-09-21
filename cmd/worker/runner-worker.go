package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc/credentials"

	pb "gitlab.com/CBCTF/bullseye-runner/proto"
	"google.golang.org/grpc"

	"gitlab.com/CBCTF/bullseye-runner/pkg/worker"
)

var (
	certFile = flag.String("cert", "", "TLS certificate")
	keyFile  = flag.String("key", "", "TLS private key")
	port     = flag.Int("port", 10080, "port to listen")
)

const Tempdir = "./tmp"

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if _, err = os.Stat(Tempdir); err != nil {
		if err2 := os.Mkdir(Tempdir, 0755); err2 != nil {
			log.Fatalf("failed to create directory: %s", Tempdir)
		}
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
	pb.RegisterRunnerServer(server, &worker.RunnerServer{})
	server.Serve(lis)
}
