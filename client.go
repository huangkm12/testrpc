package myRPC

import (
	"fmt"
	"net"
)

type Client struct {
	Conn net.Conn
}

func (client *Client) Close()  {
	client.Conn.Close()
}

func (client *Client) Call(method string,req interface{},reply interface{}) error  {

	edcode, err := GetEdCode()
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println()
	request := GetRequest(method, req)
	fmt.Println(request)
	err = request.RegisterGobArgsType()
	reqData, err := edcode.Encode(request)
	if err != nil {
		fmt.Println(err)
		return err
	}

	transfer := NewTransfer(client.Conn)
	err = transfer.write(reqData)
	if err != nil {
		fmt.Println(err)
		return err
	}

	replyData, err := transfer.read()
	if err != nil {
		fmt.Println(err)
		return err
	}

	edcode.Decode(replyData,reply)

	return nil
}
