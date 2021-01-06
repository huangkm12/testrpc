package myRPC

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
)

type Edcode interface {
	Encode(v interface{}) ([]byte,error)
	Decode([]byte,interface{}) error
}

type JsonCode int

type GobCode int

func (g GobCode) Encode(v interface{}) ([]byte, error) {
	var  buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(v)
	if err != nil {
		return nil,err
	}
	return buff.Bytes(),nil
}

func (g GobCode) Decode(data[]byte,v interface{}) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(v)
}

func (j JsonCode) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j JsonCode) Decode(data []byte,v interface{}) error {
	return json.Unmarshal(data,v)
}

func GetEdCode() (Edcode,error)  {
	EdcodeConf := new(Config).GetEdcodeConf()
	switch EdcodeConf {
	case "gob":
		return *new(GobCode),nil
	case "json":
		return *new(JsonCode),nil
	default:
		return nil,errors.New("unKnown code format ")
	}

}


