package pushTemplate

import (
	"chanPool"
	"errors"
)

type Push interface {
	Push(send Send) error
	Start()
	Stop()
}

func NewPush(chanPool chanPool.Pool, sendsnum int, resultFunc func(result Result)) Push {
	return &push{chanPool: chanPool, sends: make(chan Send, sendsnum), stopSignal: make(chan struct{}), resultFunc: resultFunc}
}

type push struct {
	stoped     bool
	stop       bool
	chanPool   chanPool.Pool
	sends      chan Send
	stopSignal chan struct{}
	resultFunc func(result Result)
}

func (this *push) Push(send Send) error {
	if this.stop {
		return errors.New("STOPED")
	}
	this.sends <- send
	return nil
}
func (this *push) Start() {
	this.chanPool.Start()
	this.chanPool.EnableWaitForAll(true)
	this.pushMessage()
}

func (this *push) Stop() {
	this.stop = true
	this.stopSignal <- struct{}{}
	<-this.stopSignal
	this.chanPool.WaitForAll()
	this.chanPool.Stop()
	close(this.sends)
	close(this.stopSignal)
	this.stoped = true
}

func (this *push) pushMessage() {
	go func() {
		ok := true
		for {
			select {
			case send, ok := <-this.sends:
				if !ok {
					break
				} else {
					this.chanPool.AddJob(func() {
						this.resultFunc(send.Send())
						return
					})
				}
			case <-this.stopSignal:
				for len(this.sends) > 0 {
					send := <-this.sends
					this.chanPool.AddJob(func() {
						this.resultFunc(send.Send())
					})
				}
				this.stopSignal <- struct{}{}
				return
			}
			if !ok {
				break
			}
		}
	}()
}
