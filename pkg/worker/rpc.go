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

	if req.PullImage {
		runner := NewRunner(ctx, req)
		err := runner.DryRun()
		if err != nil {
			return nil, err
		}

		res := &pb.RunnerResponse{}
		return res, nil
	}

	JobQueue <- struct{}{}
	runner := NewRunner(ctx, req)
	succeeded, err := runner.Run()
	<-JobQueue

	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}

	res := &pb.RunnerResponse{
		Uuid:      req.Uuid,
		Succeeded: succeeded,
	}

	return res, nil
}

func (s *RunnerServer) Info(ctx context.Context, req *pb.InfoRequest) (*pb.InfoResponse, error) {
	res := &pb.InfoResponse{
		Cpus: uint64(runtime.NumCPU()),
	}

	return res, nil
}
