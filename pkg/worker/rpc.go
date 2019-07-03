package worker

import (
	"context"
	"log"

	pb "gitlab.com/CBCTF/bullseye-runner/proto"
)

type RunnerServer struct{}

func (s *RunnerServer) Run(ctx context.Context, req *pb.RunnerRequest) (*pb.RunnerResponse, error) {
	log.Printf("received: %v", req)

	res, err := RunRequest(ctx, req)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	return res, nil
}
