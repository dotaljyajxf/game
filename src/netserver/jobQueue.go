package netserver

const MAX_WORKERS = 10
const MAX_JOB_QUEUES = 5

type callBackFunc func([]string)
type jobFunc func(string, callBackFunc)

type Job struct {
	DoJob    jobFunc
	Arg      string
	callBack callBackFunc
}

var jobQ chan Job

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
				job.DoJob(job.Arg, job.callBack)
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
}

func NewDispatcher() *Dispatcher {
	pool := make(chan chan Job, MAX_WORKERS)
	return &Dispatcher{WorkerPool: pool}
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
		case job := <-jobQ:
			go func(job Job) {
				worker := <-d.WorkerPool
				worker <- job
			}(job)
		}
	}
}
