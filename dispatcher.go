package chanPool

import (
	"errors"
)

// 协调器
type Dispatcher interface {
	Start()
	Stop()
	AddJob(job Job) error
	JobQueueLen() int
	WorkerPool() chan Worker
}

// 调度员
type dispatcher struct {
	workerPool chan Worker
	jobQueue   chan Job
	stopSignal chan struct{}
	stop       bool
	stoped     bool
}

// 创建调度器
func NewDispatcher(workerPool chan Worker, jobQueue chan Job) Dispatcher {
	return &dispatcher{workerPool: workerPool, jobQueue: jobQueue, stopSignal: make(chan struct{})}
}

// 分派工作给自由工人
func (dis *dispatcher) Start() {
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
				for len(dis.jobQueue)>0 {
					worker := <-dis.workerPool
					worker.AddJob(<-dis.jobQueue)
				}
				for len(dis.workerPool)>0 {
					worker := <-dis.workerPool
					worker.Stop()
				}
				dis.stopSignal <- struct{}{}
				return
			}
		}
	}()
}

func (dis *dispatcher) Stop() {
	dis.stop = true
	dis.stopSignal <- struct{}{}
	<-dis.stopSignal
	dis.stoped = true
}

func (dis *dispatcher) JobQueueLen() int {
	return len(dis.jobQueue)
}

func (dis *dispatcher) WorkerPool() chan Worker {
	return dis.workerPool
}

func (dis *dispatcher) AddJob(job Job) error {
	if dis.stop {
		errors.New(Stoped)
	}
	dis.jobQueue <- job
	return nil
}
