package runner

import (
	"context"
	"log"
	"time"

	pb "gitlab.com/CBCTF/bullseye-runner/proto"
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
