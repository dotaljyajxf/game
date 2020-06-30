package gateway

import (
	"sync"
)

type Call struct {
	Method  string
	Arg     []byte
	Done    chan *Call
	Ret     interface{}
	Context *TContext
}

var callPool = &sync.Pool{
	New: func() interface{} {
		return new(Call)
	},
}

var dummyCall Call

func NewCall() *Call {
	return callPool.Get().(*Call)
}

func (c *Call) Put() {
	*c = dummyCall
	callPool.Put(c)
}

var callChan chan *Call

func InitWorker(num int) {
	callChan = make(chan *Call, num)
	for i := 0; i < num; i++ {
		go runWork()
	}
}

func (this *Call) DispatchCall() {
	callChan <- this
}

func runWork() {
	for {
		select {
		case call := <-callChan:
			{
				execute(call)
				call.Done <- call
			}
		}
	}
}

func execute(call *Call) {
	HandleRequestDirect(call)
}
