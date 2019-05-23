package chanPool

import "sync"

type DefaultPool interface {
	Start() error
	Stop() error
	Close() error
	AddDefaultJob(job DefaultJob) error
	ResultChan() chan PoolError
}

type defaultPool struct {
	lock       *sync.Mutex
	dispatcher Dispatcher
	resultChan chan PoolError
}

//workerNum 工人池中的工人数量
//
//jobNum job池中的job数量
func NewDefaultPool(workerNum, jobNum, resultChanNum int) DefaultPool {
	return &defaultPool{
		lock:       new(sync.Mutex),
		dispatcher: NewDispatcher(workerNum, jobNum),
		resultChan: make(chan PoolError, resultChanNum),
	}
}

// 添加一个job到job池中
func (p *defaultPool) AddDefaultJob(defaultJob DefaultJob) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	job := func() {
		// 运行时异常处理
		err := defaultJob.Job()
		p.resultChan <- PoolError{error: err, defaultJob: defaultJob}
	}
	return p.dispatcher.AddJob(job)
}

// 停止所有进程
func (p *defaultPool) Stop() error {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.dispatcher.Stop()
}

//  清理资源
func (p *defaultPool) Close() error {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.dispatcher.Close()
}

//Start worker pool and dispatch
func (p *defaultPool) Start() (err error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	err = p.dispatcher.EnableWaitForAll(true)
	if err != nil {
		return
	}
	err = p.dispatcher.Start()
	if err != nil {
		return
	}
	return
}

func (p *defaultPool) ResultChan() chan PoolError {
	return p.resultChan
}
