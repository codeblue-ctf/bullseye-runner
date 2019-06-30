package main

import (
	"context"
	"log"

	pb "gitlab.com/CBCTF/bullseye-runner/proto"
)

type runnerServer struct{}

func (s *runnerServer) Run(ctx context.Context, req *pb.RunnerRequest) (*pb.RunnerResponse, error) {
	log.Printf("received: %v", req)

	res, err := RunRequest(ctx, req)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	return res, nil
}
