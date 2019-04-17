package pushTemplate

import (
	"github.com/JeonYang/chanPool"
)

type Push interface {
	Push(send ...Send) error
	Start()
	Stop()
}

func NewPush(chanPool chanPool.Pool, resultChanLen, resultsMaxLen int, resultFunc func(result []Result)) Push {
	return &push{chanPool: chanPool,  resultMaxLen: resultsMaxLen, resultChan: make(chan Result, resultChanLen),   resultFinishSignal: make(chan struct{}), resultFunc: resultFunc}
}

type push struct {
	chanPool           chanPool.Pool
	resultMaxLen       int
	resultChan         chan Result
	resultFinishSignal chan struct{}
	resultFunc         func(result []Result)
}

func (this *push) Push(send ...Send) error {
	for _, v := range send {
		err := this.chanPool.AddJob(func() {
			result := v.Send()
			this.resultChan <- result
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// 开启
func (this *push) Start() {
	// 开启发送消息协程
	this.chanPool.Start()
	// 开启消息反馈处理协程
	this.startResultWork()
}

// 停止
func (this *push) Stop() {
	this.chanPool.Stop()
	close(this.resultChan)
	<-this.resultFinishSignal
	close(this.resultFinishSignal)
	this.chanPool.Close()
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
