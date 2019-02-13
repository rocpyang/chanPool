package pushTemplate

import (
	"github.com/JeonYang/chanPool"
	"errors"
)
type Push interface {
	Push(send ...Send) error
	Start()
	Stop()
}

func NewPush(chanPool chanPool.Pool, sendChanLen,resultChanLen, resultsMaxLen int, resultFunc func(result []Result)) Push {
	return &push{chanPool: chanPool, sends: make(chan Send, sendChanLen), resultMaxLen: resultsMaxLen,resultChan:make(chan Result,resultChanLen), stopSignal: make(chan struct{}), resultFunc: resultFunc}
}

type push struct {
	stoped       bool
	stop         bool
	chanPool     chanPool.Pool
	sends        chan Send
	resultMaxLen int
	resultChan  chan Result
	stopSignal   chan struct{}
	resultFunc   func(result []Result)

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
func (this *push) Start() {
	this.chanPool.EnableWaitForAll(true)
	this.chanPool.Start()
	this.pushMessage()
	this.startResultWork()
}

func (this *push) Stop() {
	this.stop = true
	this.stopSignal <- struct{}{}
	<-this.stopSignal
	this.chanPool.WaitForAll()
	this.chanPool.Stop()
	close(this.sends)
	close(this.stopSignal)
	close(this.resultChan)
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
						result:=send.Send()
						//	fmt.Println("result====",result)
						this.resultChan<-result
					})
				}
			case <-this.stopSignal:
				for len(this.sends) > 0 {
					send := <-this.sends
					this.chanPool.AddJob(func() {
						result:=send.Send()
						this.resultChan<-result
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