package chanPool

import (
	"errors"
	"sync"
)

// 协调器
type Dispatcher interface {
	Start()
	Stop()
	Close() error
	AddJob(job Job) error
	WaitForAll()
	EnableWaitForAll(enable bool)
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
	stop             bool
	stoped           bool
	close            bool
	closed           bool
}

// 创建调度器
func NewDispatcher(workerNum, jobNum int) Dispatcher {
	workers := make(chan Worker, workerNum)
	jobs := make(chan Job, jobNum)
	return &dispatcher{workerPool: workers, jobQueue: jobs, stopSignal: make(chan struct{}), stopedSignal: make(chan struct{}), closeSignal: make(chan struct{}), wg: sync.WaitGroup{}, workerNum: workerNum, jobNum: jobNum}
}

// 分派工作给自由工人
func (dis *dispatcher) Start() {
	if dis.close {
		return
	}
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
	dis.stoped = false
	dis.stop = false
	dis.close = false
	dis.closed = false
	dis.workerCount = 0
}

func (dis *dispatcher) Stop() {
	if dis.close ||!dis.stoped{
		return
	}
	dis.stop = true
	dis.stopSignal <- struct{}{}
	<-dis.stopedSignal
	dis.stoped = true
}

func (dis *dispatcher) Close() error {
	if !dis.stoped {
		return errors.New("the dispatcher is start")
	}
	if dis.close {
		return nil
	}
	dis.close = true
	dis.closeSignal <- struct{}{}
	<-dis.closeSignal
	close(dis.closeSignal)
	close(dis.stopSignal)
	close(dis.jobQueue)
	close(dis.workerPool)
	dis.closed = true
	return nil
}

func (dis *dispatcher) AddJob(job Job) error {
	if dis.stop {
		return errors.New(Stoped)
	}
	if dis.close {
		return errors.New(Closed)
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
func (dis *dispatcher) WaitForAll() {
	if dis.enableWaitForAll {
		dis.wg.Wait()
	}
}

// 是否允许等待所有
func (dis *dispatcher) EnableWaitForAll(enable bool) {
	dis.enableWaitForAll = enable
}
