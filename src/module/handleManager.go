package module

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"netserver"
	"reflect"
)

type funcType struct {
	method  reflect.Method
	argType reflect.Type
	recvVal reflect.Value
	recvTyp reflect.Type
}

var funMap map[string]*funcType

func RegisterFunc(name string, moduleVar interface{}) {
	v := reflect.ValueOf(moduleVar)
	vt := v.Type()
	rcvr := v.Elem().Interface()
	if _, ok := reflect.TypeOf(rcvr).MethodByName("SetContext"); !ok {
		panic(fmt.Sprintf("module:%s need has method:%s", name, "SetContext"))
	}

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

		funMap[fullName] = &funcType{vt.Method(i), argType, reflect.ValueOf(rcvr),
			reflect.TypeOf(rcvr)}
	}
}

func HandleRequestDirect(call *netserver.Call) {

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

	rcvr := reflect.New(funcObj.recvTyp)
	setContextFunc := rcvr.MethodByName("SetContext")
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
	arrArgValues = append(arrArgValues, argv)

	ret := funcObj.method.Func.Call(arrArgValues)
	call.Ret = ret[0].Interface()
}
