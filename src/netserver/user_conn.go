package netserver

import (
	"net"
	"netserver/log"
	"pb"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"
)

type INetCodec interface {
	ReadRequest(interface{}, bool) (uint32, error)
	WriteResponse(interface{}) (uint32, error) //目前的TResponse和TRequest都很小，其实没有必要用指针，考虑以后扩展还是用了指针
	Encode(interface{}) ([]byte, error)
	JustEncode(interface{}) ([]byte, error)
	WriteQuickResponse(interface{}) (uint32, error)
	SendPing() error
	Flush() error
	Close() error
}

type UserConn struct {
	conn     *net.TCPConn
	log      *log.UserLogger
	respChan chan *pb.TResponse
	reqChan  chan *pb.TRequest
	codec    *TNetCodec
	logoff   bool //是否需要logoff，如果是自然断开的需要logoff，如果是被踢掉的就不用了
	closing  bool
	shutdown bool
	uid      uint64
	sync.Mutex
	session TSession
}

var gReponse_pool = &sync.Pool{
	New: func() interface{} {
		return new(pb.TResponse)
	},
}

var gRequest_pool = &sync.Pool{
	New: func() interface{} {
		return new(pb.TRequest)
	},
}

func NewTResponse() *pb.TResponse {
	obj := gReponse_pool.Get().(*pb.TResponse)
	obj.Reset()
	return obj
}

func NewTRequest() *pb.TRequest {
	obj := gRequest_pool.Get().(*pb.TRequest)
	obj.Reset()
	return obj
}

func FreeRequest(aReq *pb.TRequest) {
	aReq.Reset()
	gRequest_pool.Put(aReq)
}

func FreeResponse(aResp *pb.TResponse) {
	aResp.Reset()
	gReponse_pool.Put(aResp)
}

func NewUserConn(conn *net.TCPConn) *UserConn {
	log := log.NewUserLogger()
	log.AddPairInfo("client", conn.RemoteAddr().String())

	return &UserConn{
		conn:     conn,
		log:      log,
		respChan: make(chan *pb.TResponse, GlobalConfig.MaxClientReq),
		reqChan:  make(chan *pb.TRequest, GlobalConfig.MaxClientReq),
		codec:    NewNetCodec(conn),
		logoff:   false,
		closing:  false,
		shutdown: false,
		session:  TSession{},
	}
}

func (this *UserConn) Start() {
	ips := strings.Split(this.conn.RemoteAddr().String(), ":")
	this.session.Ip = ips[0]
	this.goRoutine("readRequest", this.routineReadRequest)

	this.goRoutine("handleRequest", this.HandleRequest)
	this.goRoutine("writeRequest", this.WriteResponse)
}

func (this *UserConn) HandleRequest() {
	running := true
	for running {
		select {
		case req, ok := <-this.reqChan:
			if !ok {
				this.log.Info("reqChan closed uid:%d", this.uid)
				running = false
				break
			}

			resp := NewTResponse()

			if this.uid == 0 {
				if req.GetMethod() != "user.Login" {
					this.log.Fatal("user need login")
					//TODO write a err 登录并发
					running = false
					break
				} else {
					err := this.doRequest(req.GetMethod(), req.GetArgs(), resp)
					if err != nil {
						this.log.Fatal("doRequest %s  err %s", req.GetMethod(), err)
						running = false
						break
					}
					this.respChan <- resp
				}
			} else {
				err := this.doRequest(req.GetMethod(), req.GetArgs(), resp)
				if err != nil {
					this.log.Fatal("doRequest %s  err %s", req.GetMethod(), err)
					running = false
					break
				}
				this.log.Info("doRequest %s  response %s", req.GetMethod(), resp.Ret)
				this.respChan <- resp
			}
			FreeRequest(req)
		}
	}

	this.Close(true)
}

func (this *UserConn) doRequest(aMethod string, aArgs []byte, resp *pb.TResponse) error {

	context := NewContext()

	context.StartMethod(aMethod)
	context.InitSession(this.session)

	call := NewCall()
	call.Method = aMethod
	call.Arg = aArgs
	call.Done = make(chan *Call)
	call.Context = context
	call.DispatchCall()

	<-call.Done

	if context.IsSessionChanged() {
		if this.session.Uid == 0 && context.GetSession().Uid > 0 {
			this.uid = context.GetSession().Uid
			//ret := AddUserToManager(this.id, this)
			this.log.AddPairInfo("uid", strconv.FormatUint(this.uid, 10))
		}
		if this.session.Pid != 0 && context.GetSession().Pid != this.session.Pid {
			this.log.Fatal("can't chgsession pid not equal %d != %d", this.session.Pid, context.GetSession().Pid)
		}
		if this.session.Pid == 0 && context.GetSession().Pid != 0 {
			this.log.AddPairInfo("pid", strconv.Itoa(int(context.GetSession().Pid)))
		}
		this.session = context.GetSession()
	}
	var err error
	resp.Ret, err = this.codec.Encode(call.Ret)
	if err != nil {
		return err
	}

	str := "noError"
	var errNo int32 = 0
	t := uint32(time.Now().UnixNano())
	resp.Time = &t
	resp.Method = &aMethod
	resp.Err = &errNo
	resp.ErrMsg = &str
	context.Put()
	call.Put()
	return nil
}

func (this *UserConn) WriteResponse() {
	pingInterval := time.Duration(GlobalConfig.FrontPingMs * 1000000)
	pingTimer := time.NewTimer(pingInterval)

	running := true
	for running {
		select {
		case resp, ok := <-this.respChan:
			if !ok {
				this.log.Warn("respChan closed")
				running = false
				break
			}
			this.SendRequest(resp)
		case <-pingTimer.C:
			this.SendPing()
		}
	}
	this.Close(true)
}

func (this *UserConn) SendRequest(resp *pb.TResponse) bool {
	this.Lock()
	packageLen, err := this.codec.WriteResponse(resp)
	this.Unlock()
	FreeResponse(resp)
	if err != nil {
		return false
	} else {
		this.codec.Flush()
		this.log.Info("push response to client, method:%s, size:%d", resp.GetMethod(), packageLen)
		return true
	}
}

func (this *UserConn) SendPing() error {
	this.Lock()
	defer this.Unlock()
	return this.codec.SendPing()
}

func (this *UserConn) goRoutine(name string, fun func()) {
	f := func() {
		defer func() {
			if r := recover(); r != nil {
				this.log.Fatal("client %s err %s  recover %s stack %s", this.conn.RemoteAddr().String(),
					name, r, debug.Stack())
			}
		}()
		fun()
	}
	go f()
}

func (this *UserConn) routineReadRequest() {
	this.log.Info(" start read request")

	requestNum := 0
	idleTime := time.Duration(GlobalConfig.UserIdleTimeMs * 1000000)
	for {
		if err := this.conn.SetReadDeadline(time.Now().Add(idleTime)); err != nil {
			this.log.Fatal("setReadDeadLine error %s", err)
			break
		}

		req := NewTRequest()

		packageLen, err := this.codec.ReadRequest(req, false)
		if err != nil {
			this.log.Fatal("readRequest err:%s", err)
			break
		}

		if packageLen == 0 {
			this.log.Info("read client Ping")
			break
		}

		this.reqChan <- req
		requestNum++
	}
	this.Close(true)
	this.log.Info("client %s close,requestNum %d", this.conn.RemoteAddr().String(), requestNum)
}

func (this *UserConn) Close(logoff bool) {
	this.Lock()
	if this.closing {
		this.log.Info("have closing")
		this.Unlock()
		return
	}
	this.closing = true
	this.logoff = logoff
	this.codec.Close()
	this.log.Put()
	this.Unlock()

	close(this.reqChan)
	close(this.respChan)
}
