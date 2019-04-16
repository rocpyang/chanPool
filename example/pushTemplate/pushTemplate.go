package pushTemplate

import (
	"errors"
	"github.com/JeonYang/chanPool"
)

type Push interface {
	Push(send ...Send) error
	Start()
	Stop()
}

func NewPush(chanPool chanPool.Pool, sendChanLen, resultChanLen, resultsMaxLen int, resultFunc func(result []Result)) Push {
	return &push{chanPool: chanPool, sends: make(chan Send, sendChanLen), resultMaxLen: resultsMaxLen, resultChan: make(chan Result, resultChanLen), stopSignal: make(chan struct{}), stopedSignal: make(chan struct{}), resultFinishSignal: make(chan struct{}), resultFunc: resultFunc}
}

type push struct {
	stoped             bool
	stop               bool
	chanPool           chanPool.Pool
	sends              chan Send
	resultMaxLen       int
	resultChan         chan Result
	stopSignal         chan struct{}
	stopedSignal       chan struct{}
	resultFinishSignal chan struct{}
	resultFunc         func(result []Result)
}

func (this *push) Push(send ...Send) error {
	if this.stop {
		return errors.New("STOPED")
	}
	for _, v := range send {
		this.sends <- v
	}
	return nil
}

// 开启
func (this *push) Start() {
	this.chanPool.Start()
	// 开启发送消息协程
	this.pushMessage()
	// 开启消息反馈处理协程
	this.startResultWork()
}

// 停止
func (this *push) Stop() {
	this.stop = true
	this.stopSignal <- struct{}{}
	<-this.stopedSignal
	this.chanPool.Stop()
	close(this.sends)
	close(this.stopSignal)
	close(this.stopedSignal)
	close(this.resultChan)
	<-this.resultFinishSignal
	close(this.resultFinishSignal)
	this.chanPool.Close()
	this.stoped = true
}

// 开启发送消息
func (this *push) pushMessage() {
	go func() {
		ok := true
		for {
			var send Send
			select {
			case send, ok = <-this.sends:
				if !ok {
					break
				} else {
					this.chanPool.AddJob(func() {
						result := send.Send()
						this.resultChan <- result
					})
				}
			case <-this.stopSignal:
				for len(this.sends) > 0 {
					send := <-this.sends
					this.chanPool.AddJob(func() {
						result := send.Send()
						this.resultChan <- result
					})
				}
				this.stopedSignal <- struct{}{}
				return
			}
			if !ok {
				break
			}
		}
	}()
}

// 开启反馈数据处理
func (this *push) startResultWork() {
	go func() {
		ok := true
		results := make([]Result, 0)
		for {
			var send Result
			select {
			case send, ok = <-this.resultChan:
				if !ok {
					this.resultFunc(results)
					this.resultFinishSignal <- struct{}{}
				} else {
					results = append(results, send)
					if len(results) == this.resultMaxLen {
						this.resultFunc(results)
						results = make([]Result, 0)
					}
				}
			}
			if !ok {
				break
			}
		}
	}()
}
