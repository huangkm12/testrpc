package myRPC

import (
	"net"
	"sync"
)

func NewServer() *Server  {
	return &Server{
		ServiceMap: make(map[string]map[string]*Service),
		ServerLock: sync.Mutex{},
	}
}

func NewClient(conn net.Conn) *Client  {
	return &Client{Conn:conn}
}

func Dial(network ,addr string) (*Client,error)  {
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil,err
	}

	return NewClient(conn),nil
}
