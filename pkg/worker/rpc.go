package worker

import (
	"context"
	"log"
	"runtime"

	pb "gitlab.com/CBCTF/bullseye-runner/proto"
)

type RunnerServer struct{}

var JobQueue chan struct{}

func (s *RunnerServer) Run(ctx context.Context, req *pb.RunnerRequest) (*pb.RunnerResponse, error) {
	log.Printf("received: %v", req)

	JobQueue <- struct{}{}

	res, err := RunRequest(ctx, req)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}

	<-JobQueue

	return res, nil
}

func (s *RunnerServer) Info(ctx context.Context, req *pb.InfoRequest) (*pb.InfoResponse, error) {
	res := &pb.InfoResponse{
		Cpus: uint64(runtime.NumCPU()),
	}

	return res, nil
}
