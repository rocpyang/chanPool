package pushTemplate

import (
	"chanPool"
	"fmt"
	"runtime"
	"testing"
)

func TestPush_Push(t *testing.T) {
	//debug.SetMaxThreads(5)
	runtime.GOMAXPROCS(10)
	send := NewSend(&message{}, "2010005")
	resultchan:=make(chan Result,10)
	push := NewPush(100, 100, 100, resultchan)
	push.Start()
	for i:=0;i<100000000 ;i++  {
		push.Push(send)
	}
	push.Stop()
}
//BenchmarkStringJoin1(b *testing.B)
func BenchmarkPush_Push(b *testing.B){
	send := NewSend(&message{}, "2010005")
	resultchan:=make(chan Result,10)
	push := NewPush(100, 100, 100, resultchan)
	push.Start()

	for i:=0;i<1000 ;i++  {
		push.Push(send)
	}
	push.Stop()
}
func resultFun(result Result) {
	fmt.Println(result.ErrCode()+result.ErrMsg())
}


type message struct {
	openid string
	mobile string
}

func (this *message) OpenId() string {
	return "123456"
}
func (this *message) Mobile() string {
	return "123456"
}

func TestNewPool(t *testing.T) {
	pool := chanPool.NewPool(10, 10)
	pool.Start()
	//pool.EnableWaitForAll(false)
	pool.AddJob(job_.do)
	pool.AddJob(job_.do)
	pool.AddJob(job_.do)
	pool.AddJob(job_.do)
	pool.AddJob(job_.do)
	pool.AddJob(job_.do)
	//	pool.WaitForAll()
	//time.Sleep(time.Second)

	//fmt.Println("   pool.AddJob(do)", pool.AddJob(do))
	//pool.Start()
	pool.Stop()
	//time.Sleep(time.Second * 5)
}
func do() {
	fmt.Println("=========")
}

var job_ job = job{"123", "321"}

type job struct {
	name string
	val  string
}

func (job job) do() {
	fmt.Println("name=========", job.name)
	//fmt.Println("val=========",job.val)
}

func resultChanWork(result_chan chan Result,num int)  {
	results:=make([]Result,0)
	go func() {
		ok := true
		for {
			
			select {
			case send, ok := <-result_chan:
				if !ok {
					consumptionResults(results)
					break
				} else {
					results=append(results, send)
					if len(results)== num{
						consumptionResults(results)
					}
					results=make([]Result,0)
				}
			}
			if !ok {
				break
			}
		}
	}()
	
}

func consumptionResults(results []Result)  {
	
}