package myRPC

import (
	"encoding/gob"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type Request struct {
	MethodName string
	Args interface{}
}

func GetRequest(method string,args interface{}) Request  {
	return Request{
		MethodName: method,
		Args:       args,
	}
}


// 如果是GOB编解码，则注册Args的类型，防止gob编解码错误
func (request *Request) RegisterGobArgsType() error {
	edcodeStr := new(Config).GetEdcodeConf()
	switch edcodeStr {
	case "gob":
		args := reflect.New(reflect.TypeOf(request.Args))
		if args.Kind() == reflect.Ptr {
			args = args.Elem()
		}
		gob.Register(args.Interface())
		return nil
	case "json":
		return nil
	default:
		return errors.New("Unknown edcode protocol: " + edcodeStr)
	}
}

func (request Request) MakeArgs(edcode Edcode, service Service) (reflect.Value,error){

	switch edcode.(type) {
	case GobCode:
		return reflect.ValueOf(request.Args),nil
	case JsonCode:
		reqArgs := request.Args.(map[string]interface{})
		Argsvalue := reflect.New(service.ArgType)
		err := MakeArgType(reqArgs, Argsvalue)
		if Argsvalue.Kind() == reflect.Ptr {
			Argsvalue = Argsvalue.Elem()
		}
		return Argsvalue,err
	default:
		return reflect.ValueOf(request.Args),errors.New("unknown edcode")
	}

	
}

func MakeArgType(data map[string]interface{}, obj reflect.Value) error {
	for k,v:=range data{
		err := SetField(obj,k,v)
		if err != nil {
			return err
		}
	}

	return nil
	
}

func SetField(obj reflect.Value, name string, value interface{}) error {
	structValue := obj.Elem()
	fieldByName := structValue.FieldByName(name)
	if !fieldByName.IsValid(){
		return errors.New("this field is not vailed")
	}

	if !fieldByName.CanSet(){
		return errors.New("this field can not set")
	}
	objType := fieldByName.Type()
	valueOf := reflect.ValueOf(value)
	var err error
	if objType!=valueOf.Type(){
		valueOf, err = TypeConversion(fmt.Sprintf("%d", value), fieldByName.Kind())
		if err != nil {
			return err
		}
	}

	fieldByName.Set(valueOf)
	return nil
}

// 将string类型的value值转换成reflect.Value类型
func TypeConversion(value string, ntype reflect.Kind) (reflect.Value, error) {
	switch ntype {
	case reflect.String:
		return reflect.ValueOf(value), nil
	case reflect.Int:
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	case reflect.Int8:
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int8(i)), err
	case reflect.Int16:
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int16(i)), err
	case reflect.Int32:
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int32(i)), err
	case reflect.Int64:
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	case reflect.Float32:
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(float32(i)), err
	case reflect.Float64:
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	default:
		return reflect.ValueOf(value), errors.New("unknown type：" + ntype.String())
	}
}
