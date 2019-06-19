package module

import "reflect"

type funcType struct {
	method   reflect.Method
	argType  reflect.Type
	recvType reflect.Value
}

var funMap map[string]*funcType

func RegistFunc(name string, moduleVar interface{}) {
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

func handleRequest(aMethod string, aArgs []byte, resp chan interface{}) {
	_, ok := funMap[aMethod]
	if !ok {
		panic("unkown method")
	}
}
