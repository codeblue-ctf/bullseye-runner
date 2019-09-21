package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

var (
	flagPath = flag.String("flagpath", "/flag", "path to flag")
	port     = flag.Int("port", 1337, "port")
)

func main() {
	flag.Parse()

	file, err := os.Create(*flagPath)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	defer listen.Close()

	if err != nil {
		log.Fatalf("failed to open tcp: %v", err)
	}

	conn, err := listen.Accept()

	if err != nil {
		log.Fatalf("failed to accept request: %v", err)
	}

	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
	if err != nil {
		log.Fatalf("failed to read: %v", err)
	}

	_, err = file.Write(buf)
	if err != nil {
		log.Fatalf("failed to write: %v", err)
	}
}
