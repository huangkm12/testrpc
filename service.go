package myRPC

import (
	"encoding/gob"
	"errors"
	"reflect"
)

type Service struct {
	Method reflect.Method
	ArgType reflect.Type
	ReplyType reflect.Type
}
// gob编码需要注册ARGS的类型，防止编解码错误
func (service *Service) RegisterGobArgsType() error {
	edcodeConf := new(Config).GetEdcodeConf()
	switch edcodeConf {
	case "gob":
		value := reflect.New(service.ArgType)
		if value.Kind()==reflect.Ptr{
			value=value.Elem()
		}
		gob.Register(value.Interface())
		return nil
	case "json":
		return nil
	default:
		return errors.New("Unknown edcode protocl:"+edcodeConf)
	}
}
