package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc/credentials"

	pb "gitlab.com/CBCTF/bullseye-runner/proto"
	"google.golang.org/grpc"
)

var (
	caFile = flag.String("ca", "", "CA root cert file")
	host   = flag.String("host", "127.0.0.1:10080", "Server address")
)

func sendRequest(client pb.RunnerClient, req *pb.RunnerRequest) (*pb.RunnerResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := client.Run(ctx, req)
	if err != nil {
		log.Fatalf("%v.Run(_) = _, %v", client, err)
	}
	log.Printf("%v", res)

	return nil, nil
}

func main() {
	flag.Parse()
	var opts []grpc.DialOption

	if *caFile != "" {
		creds, err := credentials.NewClientTLSFromFile(*caFile, *host)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(*host, opts...)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	defer conn.Close()
	client := pb.NewRunnerClient(conn)

	req := pb.RunnerRequest{
		Uuid:                "hoge",
		Timeout:             1000,
		DockerComposeYml:    `hoge`,
		DockerRegistryToken: "test",
		FlagTemplate:        "CBCTF{hoge}",
		CallbackUrl:         "http://hogehoge",
		CallbackAuthToken:   "test",
	}

	sendRequest(client, &req)

}
