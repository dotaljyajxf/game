package netserver

import (
	"sync"
)

const MAX_WORKERS = 20
const MAX_JOB_QUEUES = 30

type callBackFunc func([]string)
type jobFunc func(aMethod string, aArgs []byte,resp chan interface{})

type Job struct {
	DoJob  jobFunc
	Method string
	Arg    []byte
	resp   chan interface{}
}

var jobPool = &sync.Pool{
	New: func() interface{} {
		return new(Job)
	},
}

func NewJob() *Job {
	return jobPool.Get().(*Job)
}

var jobDummy Job
var jobDispacher *Dispatcher

func (j *Job) Release() {
	*j = jobDummy
	jobPool.Put(j)
}

type Worker struct {
	WorkQueue chan chan Job
	JobQueue  chan Job
	Quit      chan bool
}

func NewWorker(workPool chan chan Job) Worker {
	return Worker{
		WorkQueue: workPool,
		JobQueue:  make(chan Job, MAX_JOB_QUEUES),
		Quit:      make(chan bool),
	}
}

func (w Worker) Start() {
	go func() {
		for {

			w.WorkQueue <- w.JobQueue

			select {
			case job := <-w.JobQueue:
				job.DoJob(job.Method,job.Arg,job.resp)
			case <-w.Quit:
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.Quit <- true
	}()
}

type Dispatcher struct {
	WorkerPool chan chan Job
	JobChan    chan Job
}

func NewDispatcher() *Dispatcher {
	pool := make(chan chan Job, MAX_WORKERS)
	return &Dispatcher{WorkerPool: pool, JobChan: make(chan Job, MAX_JOB_QUEUES*MAX_WORKERS)}
}

func (d *Dispatcher) AddJob(j Job) {
	d.JobChan <- j
}

func (d *Dispatcher) Run() {
	for i := 0; i < MAX_WORKERS; i++ {
		worker := NewWorker(d.WorkerPool)
		worker.Start()
	}
	go d.Dispatcher()
}

func (d *Dispatcher) Dispatcher() {
	for {
		select {
		case job := <-d.JobChan:
			go func(job Job) {
				worker := <-d.WorkerPool
				worker <- job
			}(job)
		}
	}
}
