package chanPool

// 工人
type Worker interface {
	Start()
	Stop()
	AddJob(job Job)
}

type worker struct {
	workerPool chan Worker
	jobQueue   chan Job
	stopSignal chan struct{}
	stoped     bool
	stop       bool
}

func NewWorker(workerPool chan Worker) *worker {
	return &worker{
		workerPool: workerPool,
		jobQueue:   make(chan Job),
		stopSignal: make(chan struct{}),
	}
}

// 工人开始工作
func (w *worker) Start() {
	go func() {
		for {
			// 将本身注册给相应的工人池
			w.workerPool <- w
			select {
			// 监听工人的工作通道
			case job := <-w.jobQueue:
				job()
				//监听工人停止信号
			case <-w.stopSignal:
				w.stopSignal <- struct{}{}
				return
			}

		}
	}()

}

// 工人开始工作
func (w *worker) Stop() {
	w.stopSignal <- struct{}{}
	<-w.stopSignal
	close(w.stopSignal)
	close(w.jobQueue)
}

// 工人开始工作
func (w *worker) AddJob(job Job) {
	w.jobQueue <- job
}

