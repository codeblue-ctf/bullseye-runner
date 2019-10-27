package worker

import (
	"context"
	"io/ioutil"
	"log"
	"os"
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

	var x11cap []byte
	if runner.x11capturing {
		x11cap, err = ioutil.ReadFile(runner.x11capPath)
		if err != nil {
			return nil, err
		}
		if err := os.Remove(runner.x11capPath); err != nil {
			return nil, err
		}
	}

	res := &pb.RunnerResponse{
		Uuid:      req.Uuid,
		Succeeded: succeeded,
		X11Cap:    x11cap,
	}

	return res, nil
}

func (s *RunnerServer) Info(ctx context.Context, req *pb.InfoRequest) (*pb.InfoResponse, error) {
	res := &pb.InfoResponse{
		Cpus: uint64(runtime.NumCPU()),
	}

	return res, nil
}
