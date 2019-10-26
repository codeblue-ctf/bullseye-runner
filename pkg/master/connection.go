package master

import (
	"sync"

	"google.golang.org/grpc"
)

var (
	cpmut sync.Mutex
)

type ConnPool struct {
	m sync.Map
}

func NewConnPool() *ConnPool {
	return &ConnPool{}
}

func (c *ConnPool) HasHost(host string) bool {
	_, ok := c.m.Load(host)
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

	c.m.Store(host, conn)
	return nil
}

func (c *ConnPool) GetConn(host string) (*grpc.ClientConn, error) {
	conn, ok := c.m.Load(host)

	if !ok {
		err := c.AddHost(host)
		if err != nil {
			return nil, err
		}
		return c.GetConn(host)
	}

	return conn.(*grpc.ClientConn), nil
}
