package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net"
)

func main() {
	flag.Parse()

	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer listen.Close()

	conn, err := listen.Accept()
	if err != nil {
		log.Fatalf("failed to acccept: %v", err)
	}
	defer conn.Close()

	flagBytes, err := ioutil.ReadFile("/flag")
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	_, err = conn.Write(flagBytes)
	if err != nil {
		log.Fatalf("failed to write socket: %v", err)
	}
}
