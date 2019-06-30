package main

import (
	"flag"

	"google.golang.org/grpc"
)

func main() {
	flag.Parse()

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())
}
