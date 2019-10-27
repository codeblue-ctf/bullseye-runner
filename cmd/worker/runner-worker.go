package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"

	"google.golang.org/grpc/credentials"

	pb "gitlab.com/CBCTF/bullseye-runner/proto"
	"google.golang.org/grpc"

	"gitlab.com/CBCTF/bullseye-runner/pkg/worker"
)

var (
	certFile   = flag.String("cert", "", "TLS certificate")
	keyFile    = flag.String("key", "", "TLS private key")
	port       = flag.Int("port", 10080, "port to listen")
	xvfbpath   = flag.String("xvfbpath", "/usr/bin/Xvfb", "path to Xvfb binary")
	ffmpegpath = flag.String("ffmpegpath", "/usr/bin/ffmpeg", "path to ffmpeg")
)

func initJobQueue() {
	cpus := runtime.NumCPU()
	worker.JobQueue = make(chan struct{}, cpus*2)
}

func checkXvfb() {
	_, err := os.Stat(*xvfbpath)
	if os.IsNotExist(err) {
		panic("Xvfb does not exist")
	}
	worker.XvfbPath = *xvfbpath
}

func checkFfmpeg() {
	_, err := os.Stat(*ffmpegpath)
	if os.IsNotExist(err) {
		panic("ffmpeg does not exist")
	}
	worker.FFmpegPath = *ffmpegpath
}

func main() {
	flag.Parse()

	checkXvfb()
	checkFfmpeg()
	worker.InitXvfb()

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

	initJobQueue()

	server := grpc.NewServer(opts...)
	pb.RegisterRunnerServer(server, &worker.RunnerServer{})
	server.Serve(lis)
}
