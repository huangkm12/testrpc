package myRPC

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"strings"
	"sync"
)

type Server struct {
	ServiceMap map[string]map[string]*Service
	ServerLock sync.Mutex
	ServerType reflect.Type
}

func (server *Server) Register(obj interface{}) error  {
	// 先进行加锁
	server.ServerLock.Lock()
	defer server.ServerLock.Unlock()
	// 根据obj来进行注册
	objType := reflect.TypeOf(obj)
	value := reflect.ValueOf(obj)

	serivceName := reflect.Indirect(value).Type().Name()

	if _,ok:=server.ServiceMap[serivceName];ok{
		return errors.New("service has register,please repeat do this")
	}

	services := make(map[string]*Service,objType.NumMethod())
	// 对结构体方法进行注册
	for i:=0;i<objType.NumMethod();i++{
		method := objType.Method(i)
		methodName := method.Name
		methodType := method.Type

		var service Service
		service.Method=method
		service.ArgType=methodType.In(1)
		service.ReplyType=methodType.In(2)

		services[methodName]=&service
		err := service.RegisterGobArgsType()
		if err != nil {
			return nil
		}
	}

	server.ServiceMap[serivceName]=services
	server.ServerType=reflect.TypeOf(obj)
	return nil

}

func (server *Server) Server(network string,addr string) error  {
	listener, err := net.Listen(network, addr)
	if err != nil {
		return err
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go server.ConnHandle(conn)
	}
}

func (server *Server)ConnHandle(conn net.Conn) {
	transfer := Transfer{conn: conn}
	for {
		data, err := transfer.read()
		if err != nil {
			fmt.Println("read err:",err)
			return
		}
		edcode, err := GetEdCode()
		if err != nil {
			fmt.Println("get edcode failed:",err)
			return
		}
		var requet Request

		err = edcode.Decode(data, &requet)
		if err!=nil{
			fmt.Println("decode failed:",err)
			return
		}

		stringSlice := strings.Split(requet.MethodName, ".")
		if len(stringSlice)!=2{
			fmt.Println("service name parsed failed")
			return
		}

		service := server.ServiceMap[stringSlice[0]][stringSlice[1]]

		argsv,err := requet.MakeArgs(edcode,*service)

		replyValue := reflect.New(service.ReplyType.Elem())
		fmt.Println(argsv)
		CallFunc := service.Method.Func

		out := CallFunc.Call([]reflect.Value{reflect.New(server.ServerType.Elem()), argsv, replyValue})
		if out[0].Interface()!=nil{
			return
		}
		fmt.Println(replyValue.Elem().Interface())
		resp, err := edcode.Encode(replyValue.Elem().Interface())
		if err != nil {
			return
		}

		err = transfer.write(resp)
		if err != nil {
			return
		}

	}

}
