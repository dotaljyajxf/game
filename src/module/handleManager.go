package module

import (
	"github.com/golang/protobuf/proto"
	"reflect"
)

type funcType struct {
	method   reflect.Method
	argType  reflect.Type
	recvType reflect.Value
}

var funMap map[string]*funcType

func RegisterFunc(name string, moduleVar interface{}) {
	v := reflect.ValueOf(moduleVar)
	vt := v.Type()

	funcNum := v.NumMethod()

	for i := 0; i < funcNum; i++ {
		funcName := vt.Method(i).Name
		fullName := name + "." + funcName
		_, ok := funMap[fullName]
		if ok {
			return
		}

		fType := vt.Method(i).Type
		argNum := fType.NumIn()
		if argNum > 1 {
			panic("rpc function must one arg")
		}

		argType := fType.In(0)

		funMap[fullName] = &funcType{vt.Method(i), argType, v}
	}
}

func HandleRequest(aMethod string, aArgs []byte, resp chan interface{}) {
	_, ok := funMap[aMethod]
	if !ok {
		panic("unkown method")
	}
	funcObj := funMap[aMethod]
	arg := reflect.New(funcObj.argType)

	err := proto.Unmarshal(aArgs, arg.Interface().(proto.Message))
	if err != nil {
		panic("args Unmarshal err")
	}

	arrArgValues := make([]reflect.Value, 0)
	arrArgValues = append(arrArgValues, arg)

	ret := funcObj.method.Func.Call(arrArgValues)

	resp <- ret[0].Interface()
}

func HandleRequestDirect(aMethod string, aArgs []byte) interface{} {
	_, ok := funMap[aMethod]
	if !ok {
		panic("unkown method")
	}
	funcObj := funMap[aMethod]
	arg := reflect.New(funcObj.argType)

	err := proto.Unmarshal(aArgs, arg.Interface().(proto.Message))
	if err != nil {
		panic("args Unmarshal err")
	}

	arrArgValues := make([]reflect.Value, 0)
	arrArgValues = append(arrArgValues, arg)

	ret := funcObj.method.Func.Call(arrArgValues)

	return ret[0].Interface()
}
