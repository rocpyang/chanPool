package chanPool

import (
	"sync"
	"errors"
)

type Pool interface {
	Start()
	Stop()
	AddJob(job Job) error
	WaitForAll()
	EnableWaitForAll(enable bool)
}

type pool struct {
	dispatcher       Dispatcher
	wg               sync.WaitGroup
	enableWaitForAll bool // 启用所有等待
	workerNum        int  // 工人总数
	jobNum           int  // 工作数
	workerCount      int  // 正在工作工人的数量
	stoped           bool
	stop             bool
}

//workerNum 工人池中的工人数量
//
//jobNum job池中的job数量
func NewPool(workerNum, jobNum int) Pool {
	workers := make(chan Worker, workerNum)
	jobs := make(chan Job, jobNum)
	return &pool{
		dispatcher:       NewDispatcher(workers, jobs),
		enableWaitForAll: false,
		workerNum:        workerNum,
		jobNum:           jobNum,
	}
}

// 添加一个job到job池中
func (p *pool) AddJob(job Job) error {
	if p.stop {
		return errors.New(Stoped)
	}
	if p.enableWaitForAll {
		p.wg.Add(1)
	}
	err := p.dispatcher.AddJob(func() {
		job()
		if p.enableWaitForAll {
			p.wg.Done()
		}
	})
	if err != nil {
		return err
	}
	if p.dispatcher.JobQueueLen() > 0 {
		if p.workerCount < p.workerNum {
			worker := NewWorker(p.dispatcher.WorkerPool())
			worker.Start()
			p.workerCount++
		}
	}
	return nil
}

// 等待所有协程操作完成
func (p *pool) WaitForAll() {
	if p.enableWaitForAll {
		p.wg.Wait()
	}
}

// 停止所有进程
func (p *pool) Stop() {
	p.stop = true
	p.dispatcher.Stop()
	p.stoped = true
	p.workerCount = 0
}

// 是否允许等待所有
func (p *pool) EnableWaitForAll(enable bool) {
	p.enableWaitForAll = enable
}

//Start worker pool and dispatch
func (p *pool) Start() {
	p.dispatcher.Start()
	p.stoped = false
	p.stop = false
}

