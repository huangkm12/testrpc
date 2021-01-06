package myRPC

import (
	"bytes"
	"net"
)
// 每次读取50个字节
const EachReadBytes  = 50
// 传输层
type Transfer struct {
	conn net.Conn
}

func NewTransfer(conn net.Conn) *Transfer  {
	return &Transfer{conn:conn}
}

func (trans *Transfer) read()  ([]byte,error){
	var buff bytes.Buffer
	bytesData := make([]byte,EachReadBytes)
	for{

		n, err := trans.conn.Read(bytesData)
		if err != nil {
			return nil,err
		}
		buff.Write(bytesData[0:n])
		if n < EachReadBytes{
			break
		}
	}
	return buff.Bytes(),nil
}

func (trans Transfer) write(data []byte) error {
	_,err:= trans.conn.Write(data)
	return err
}
