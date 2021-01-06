package main

import (
	myRPC "github.com/huangkm12/testrpc"
	"log"
)

type Args struct {
	A, B int
}

func main() {

	client, err := myRPC.Dial("tcp", "127.0.0.1:1234")

	args := Args{
		A: 5,
		B: 6,
	}
	defer client.Close()
	//defer client.Conn.Close()
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	log.Println(reply)


}
