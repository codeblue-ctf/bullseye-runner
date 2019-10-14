package master

import (
	"fmt"
	"sync"

	"google.golang.org/grpc"
)

var (
	pmut sync.Mutex
)

type ConnPool struct {
	connections map[string]*grpc.ClientConn
}

func (c *ConnPool) HasHost(host string) bool {
	pmut.Lock()
	defer pmut.Unlock()

	_, ok := c.connections[host]
	return ok
}

func (c *ConnPool) AddHost(host string) error {
	if c.HasHost(host) == true {
		return nil
	}

	pmut.Lock()
	defer pmut.Unlock()

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
	pmut.Lock()
	defer pmut.Unlock()

	conn, ok := c.connections[host]

	if !ok {
		return nil, fmt.Errorf("no such connection: %s", host)
	}

	return conn, nil
}
