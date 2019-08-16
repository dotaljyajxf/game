package netserver

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"reflect"
)

type funcType struct {
	method  reflect.Method
	argType reflect.Type
	recvVal reflect.Value
	recvTyp reflect.Type
}

var funMap = make(map[string]*funcType, 300)

func RegisterFunc(name string, moduleVar interface{}) {
	v := reflect.ValueOf(moduleVar)
	vt := v.Type()
	rcvr := v.Elem().Interface()
	if _, ok := vt.MethodByName("SetContext"); !ok {
		panic(fmt.Sprintf("module:%s need has method:%s", name, "SetContext"))
	}

	funcNum := v.NumMethod()

	for i := 0; i < funcNum; i++ {
		funcName := vt.Method(i).Name

		if funcName == "SetContext" || funcName == "GetContext" {
			continue
		}
		fullName := name + "." + funcName
		_, ok := funMap[fullName]
		if ok {
			return
		}
		fType := vt.Method(i).Type
		argNum := fType.NumIn()

		if argNum != 2 {
			panic("rpc function must one arg")
		}

		argType := fType.In(1)

		funMap[fullName] = &funcType{vt.Method(i), argType, v,
			reflect.TypeOf(rcvr)}
	}
}

func HandleRequestDirect(call *Call) {

	context := call.Context
	aMethod := call.Method
	aArgs := call.Arg
	aLog := context.GetLogger()

	_, ok := funMap[aMethod]
	if !ok {
		aLog.Fatal("not found method!")
		return
	}
	funcObj := funMap[aMethod]

	var argv reflect.Value
	isVal := false
	if funcObj.argType.Kind() == reflect.Ptr {
		argv = reflect.New(funcObj.argType.Elem())
	} else {
		isVal = true
		argv = reflect.New(funcObj.argType)
	}

	setContextFunc := funcObj.recvVal.MethodByName("SetContext")
	setContextFunc.Call([]reflect.Value{reflect.ValueOf(context)})

	err := proto.Unmarshal(aArgs, argv.Interface().(proto.Message))
	if err != nil {
		aLog.Fatal("arg type error")
		return
	}

	if isVal {
		argv = argv.Elem()
	}

	arrArgValues := make([]reflect.Value, 0)
	arrArgValues = append(arrArgValues, funcObj.recvVal)
	arrArgValues = append(arrArgValues, argv)

	ret := funcObj.method.Func.Call(arrArgValues)
	call.Ret = ret[0].Interface()
}
