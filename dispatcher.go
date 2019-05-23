package chanPool

import (
	"sync"
)

// 协调器
type Dispatcher interface {
	Start() error
	Stop() error
	Close() error
	AddJob(job Job) error
	WaitForAll() error
	EnableWaitForAll(enable bool) error
}

// 调度员
type dispatcher struct {
	workerPool       chan Worker
	jobQueue         chan Job
	stopSignal       chan struct{}
	stopedSignal     chan struct{}
	closeSignal      chan struct{}
	wg               sync.WaitGroup
	enableWaitForAll bool // 启用所有等待
	workerNum        int  // 工人总数
	jobNum           int  // 工作数
	workerCount      int  // 正在工作工人的数量
	start            bool
	close            bool
}

// 创建调度器
func NewDispatcher(workerNum, jobNum int) Dispatcher {
	workers := make(chan Worker, workerNum)
	jobs := make(chan Job, jobNum)
	return &dispatcher{workerPool: workers, jobQueue: jobs, stopSignal: make(chan struct{}), stopedSignal: make(chan struct{}), closeSignal: make(chan struct{}), wg: sync.WaitGroup{}, workerNum: workerNum, jobNum: jobNum}
}

// 分派工作给自由工人
func (dis *dispatcher) Start() error {
	if dis.start {
		return START
	}
	if dis.close {
		return CLOSED
	}
	dis.start = true
	dis.close = false
	dis.workerCount = 0
	go func() {
		for {
			select {
			// 监听调度器的工作通道
			case job := <-dis.jobQueue:
				worker := <-dis.workerPool
				worker.AddJob(job)
				// 监听调度器的停止信号
			case <-dis.stopSignal:
				for len(dis.jobQueue) > 0 {
					worker := <-dis.workerPool
					worker.AddJob(<-dis.jobQueue)
				}
				// 停止所有的工人
				for dis.workerCount > 0 {
					worker := <-dis.workerPool
					worker.Stop()
					dis.workerCount = dis.workerCount - 1
				}
				dis.stopedSignal <- struct{}{}
				// 监听资源清除信号
			case <-dis.closeSignal:

				dis.closeSignal <- struct{}{}
				return
			}
		}
	}()
	return nil
}

func (dis *dispatcher) Stop() error {
	if dis.close {
		return CLOSED
	}
	if !dis.start {
		return nil
	}
	dis.stopSignal <- struct{}{}
	<-dis.stopedSignal
	dis.start = false
	return nil
}

func (dis *dispatcher) Close() error {
	if dis.start {
		return START
	}
	if dis.close {
		return nil
	}
	dis.closeSignal <- struct{}{}
	<-dis.closeSignal
	close(dis.closeSignal)
	close(dis.stopSignal)
	close(dis.jobQueue)
	close(dis.workerPool)
	dis.close = true
	return nil
}

func (dis *dispatcher) AddJob(job Job) error {
	if !dis.start {
		return STOPED
	}
	if dis.close {
		return CLOSED
	}
	if dis.enableWaitForAll {
		dis.wg.Add(1)
	}
	dis.jobQueue <- func() {
		job()
		if dis.enableWaitForAll {
			dis.wg.Done()
		}
	}
	if len(dis.jobQueue) > 0 || dis.workerCount == 0 {
		if dis.workerCount < dis.workerNum {
			worker := NewWorker(dis.workerPool)
			worker.Start()
			dis.workerCount++
		}
	}
	return nil
}

// 等待所有协程操作完成
func (dis *dispatcher) WaitForAll() error {
	if !dis.start {
		return STOPED
	}
	if dis.close {
		return CLOSED
	}
	if dis.enableWaitForAll {
		dis.wg.Wait()
	}
	return nil
}

// 是否允许等待所有
func (dis *dispatcher) EnableWaitForAll(enable bool) error {
	if dis.start {
		return START
	}
	dis.enableWaitForAll = enable
	return nil
}
