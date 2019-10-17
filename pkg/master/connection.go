package master

import (
	"fmt"
	"sync"

	"google.golang.org/grpc"
)

var (
	cpmut sync.Mutex
)

type ConnPool struct {
	connections map[string]*grpc.ClientConn
}

func NewConnPool() ConnPool {
	pool := ConnPool{}
	pool.connections = make(map[string]*grpc.ClientConn)
	return pool
}

func (c *ConnPool) HasHost(host string) bool {
	cpmut.Lock()
	defer cpmut.Unlock()

	_, ok := c.connections[host]
	return ok
}

func (c *ConnPool) AddHost(host string) error {
	if c.HasHost(host) == true {
		return nil
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(host, opts...)
	if err != nil {
		return err
	}

	cpmut.Lock()
	defer cpmut.Unlock()

	c.connections[host] = conn
	return nil
}

func (c *ConnPool) GetConn(host string) (*grpc.ClientConn, error) {
	cpmut.Lock()
	defer cpmut.Unlock()

	conn, ok := c.connections[host]

	if !ok {
		return nil, fmt.Errorf("no such connection: %s", host)
	}

	return conn, nil
}
