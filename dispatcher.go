package chanPool

import (
	"errors"
)

// 协调器
type Dispatcher interface {
	Start()
	Close() error
	Stop()
	AddJob(job Job) error
	JobQueueLen() int
	WorkerPool() chan Worker
}

// 调度员
type dispatcher struct {
	workerPool  chan Worker
	jobQueue    chan Job
	stopSignal  chan struct{}
	closeSignal chan struct{}
	stop        bool
	stoped      bool
	close       bool
	closed      bool
}

// 创建调度器
func NewDispatcher(workerPool chan Worker, jobQueue chan Job) Dispatcher {
	return &dispatcher{workerPool: workerPool, jobQueue: jobQueue, stopSignal: make(chan struct{}), closeSignal: make(chan struct{})}
}

// 分派工作给自由工人
func (dis *dispatcher) Start() {
	if dis.close {
		return
	}
	dis.stoped = false
	dis.stop = false
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
				for len(dis.workerPool) > 0 {
					worker := <-dis.workerPool
					worker.Stop()
				}
				dis.stopSignal <- struct{}{}
				// 监听资源清除信号
			case <-dis.closeSignal:
				dis.closeSignal <- struct{}{}
				return
			}
		}
	}()
}

func (dis *dispatcher) Stop() {
	if dis.close {
		return
	}
	dis.stop = true
	dis.stopSignal <- struct{}{}
	<-dis.stopSignal
	dis.stoped = true
}

func (dis *dispatcher) Close() error {
	if !dis.stoped {
		return errors.New("the dispatcher is start")
	}
	dis.close = true
	dis.closeSignal <- struct{}{}
	<-dis.closeSignal
	close(dis.closeSignal)
	close(dis.stopSignal)
	dis.closed = true
	return nil
}

func (dis *dispatcher) JobQueueLen() int {
	return len(dis.jobQueue)
}

func (dis *dispatcher) WorkerPool() chan Worker {
	return dis.workerPool
}

func (dis *dispatcher) AddJob(job Job) error {
	if dis.stop {
		return errors.New(Stoped)
	}
	dis.jobQueue <- job
	return nil
}
