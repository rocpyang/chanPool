package chanPool

import "sync"

type Pool interface {
	Start() error
	Stop() error
	Close() error
	AddJob(job Job) error
}

type pool struct {
	lock       *sync.Mutex
	dispatcher Dispatcher
}

//workerNum 工人池中的工人数量
//
//jobNum job池中的job数量
func NewPool(workerNum, jobNum int) Pool {
	return &pool{
		lock:       new(sync.Mutex),
		dispatcher: NewDispatcher(workerNum, jobNum),
	}
}

// 添加一个job到job池中
func (p *pool) AddJob(job Job) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.dispatcher.AddJob(job)
}

// 停止所有进程
func (p *pool) Stop() error {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.dispatcher.Stop()
}

//  清理资源
func (p *pool) Close() error {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.dispatcher.Close()
}

//Start worker pool and dispatch
func (p *pool) Start() (err error) {
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
