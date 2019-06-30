package main

import (
	"fmt"

	"google.golang.org/grpc"
)

type ConnPool struct {
	connections map[string]*grpc.ClientConn
}

func (c *ConnPool) AddHost(host string) error {
	if _, ok := c.connections[host]; ok {
		return nil
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(host, opts...)
	if err != nil {
		return err
	}

	c.connections[host] = conn
	return nil
}

func (c *ConnPool) GetConn(host string) (*grpc.ClientConn, error) {
	conn, ok := c.connections[host]

	if !ok {
		return nil, fmt.Errorf("no such connection: %s", host)
	}

	return conn, nil
}
