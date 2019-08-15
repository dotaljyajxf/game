package codeTest

import (
	"golang.org/x/tools/go/ssa/interp/testdata/src/fmt"
	"net"
	"pb"
)

var wanAddr = "192.168.1.185:2000"

type UserClient struct {
}

func iniClient() *TNetCode {

	tcpAddr, err := net.ResolveTCPAddr("tcp4", wanAddr)
	if err != nil {
		panic("resolve addr error!")
	}

	conn, err := net.DialTCP("tcp4", nil, tcpAddr)
	if err != nil {
		panic("dial to tcpAddr.IP.String() error")
	}

	return NewTNetCode(conn)
}

func (this *UserClient) SendReq(aMethod string, arg interface{}, resp interface{}) bool {
	client := iniClient()

	pbCodec := &TPbCode{}

	argByte, err := pbCodec.EnCode(arg)
	if err != nil {
		panic("arg encode error")
	}

	token := "1"
	req := &pb.TRequest{
		Method:          &aMethod,
		Args:            argByte,
		CallbackHandler: nil,
		Token:           &token,
	}

	reqByte, err := pbCodec.EnCode(req)
	if err != nil {
		panic("req encode error")
	}

	err = client.WriteRequest(reqByte, false)
	if err != nil {
		panic(err)
	}

	for {
		resByte, err := client.ReadResponse()
		if err != nil {
			panic("read response err")
		}

		if resByte == nil {
			continue
		}

		respRet := &pb.TResponse{}
		err = pbCodec.DeCode(resByte, respRet)

		if respRet.GetErr() != 0 {
			fmt.Println("run method[%s] err[%d]:errmsg[%s]", respRet.GetMethod(), respRet.GetErr(), respRet.GetErrMsg())
			return false
		}

		if *respRet.Method == *req.Method {
			pbCodec.DeCode(respRet.Ret, resp)
			break
		}
	}
	return true
}
