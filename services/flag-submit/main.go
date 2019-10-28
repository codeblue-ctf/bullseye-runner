package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

var (
	flagPath = flag.String("flagpath", "/flag", "path to flag")
	port     = flag.Int("port", 1337, "port")
)

func main() {
	flag.Parse()

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
	if _, err := conn.Read(buf); err != nil {
		log.Fatalf("failed to read: %v (buf: %+v)", err, buf)
	}

	if err := ioutil.WriteFile(*flagPath, buf, 0644); err != nil {
		log.Fatalf("failed to write file: %+v", err)
	}
}
