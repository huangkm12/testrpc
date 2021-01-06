package main

import (
	"github.com/huangkm12/testrpc"
	"net/rpc"
)


type Args struct {
	A, B int
}

type Arith int

func (t *Arith) Multiply(args Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func main() {
	server := myRPC.NewServer()
	server.Register(new(Arith))
	rpc.NewClient()
	server.Server("tcp",":1234")
}
