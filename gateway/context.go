package gateway

import (
	"game/data"
	"game/gateway/log"
	"github.com/go-xorm/xorm"
	"strconv"
	"sync"
	"time"
)

type TSession struct {
	Uid      uint64 //TODO: 一个角色id，一个账号id。是否需要调整一个名字
	Pid      uint32
	ServerId int32
	Ip       string
}

var dummyContext TContext

var contextPool = &sync.Pool{
	New: func() interface{} {
		return new(TContext)
	},
}

func NewContext() *TContext {
	return contextPool.Get().(*TContext)
}

func (c *TContext) Put() {
	c.logger.Put()
	*c = dummyContext
	contextPool.Put(c)
}

type TContext struct {
	logger       *log.UserLogger
	moduleMethod string
	reqStartTime time.Time
	reqEndTime   time.Time
	session      TSession
	isSessionChg bool
	mData        *xorm.Engine
}

func (this *TContext) SetSession(session TSession) {
	this.session = session
}

func (this *TContext) StartMethod(method string) {
	this.moduleMethod = method
	this.reqStartTime = time.Now()
	this.isSessionChg = false
	this.logger = log.NewUserLogger()
	this.mData = data.NewEngine()
}

func (this *TContext) DB() *xorm.Engine {
	return this.mData
}

func (this *TContext) GetLogger() *log.UserLogger {
	return this.logger
}

func (this *TContext) GetModuleMethod() string {
	return this.moduleMethod
}

func (this *TContext) IsSessionChanged() bool {
	return this.isSessionChg
}

func (this *TContext) GetSession() TSession {
	return this.session
}

func (this *TContext) InitSession(aSession TSession) {
	this.session = aSession
	if aSession.Uid != 0 {
		this.logger.AddPairInfo("uid", strconv.Itoa(int(this.session.Uid)))
	}
	if aSession.Pid != 0 {
		this.logger.AddPairInfo("pid", strconv.Itoa(int(this.session.Pid)))
	}
}

func (this *TContext) Now() time.Time {
	return time.Unix(this.Unix(), int64(this.reqStartTime.Nanosecond()))
}

func (this *TContext) Unix() int64 {
	return this.reqStartTime.Unix()
}

func (this *TContext) UnixNano() int64 {
	return this.reqStartTime.UnixNano()
}

func (this *TContext) GetPid() uint32 {
	return this.session.Pid
}

func (this *TContext) SetPid(aPid uint32) {
	if this.session.Pid == aPid {
		return
	}
	this.session.Pid = aPid
	this.isSessionChg = true
	this.logger.AddPairInfo("pid", strconv.Itoa(int(this.session.Pid)))
}

func (this *TContext) GetUid() uint64 {
	return this.session.Uid
}

func (this *TContext) SetUid(aUid uint64) {
	if this.session.Uid == aUid {
		return
	}
	this.session.Uid = aUid
	this.isSessionChg = true
	this.logger.AddPairInfo("uid", strconv.Itoa(int(this.session.Uid)))
}

func (this *TContext) GetLoginIp() string {
	return this.session.Ip
}

func (this *TContext) ClearLoginIp() {
	if len(this.session.Ip) <= 0 {
		return
	}
	this.session.Ip = ""
	this.isSessionChg = true
}
