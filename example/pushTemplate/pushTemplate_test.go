package pushTemplate

import (
	"fmt"
	"github.com/JeonYang/chanPool"
	"runtime"
	"testing"
)

func TestPush_Push(t *testing.T) {
	//debug.SetMaxThreads(5)
	runtime.GOMAXPROCS(10)
	send := NewSend(&message{}, "2010005")
	pool:=chanPool.NewPool(100,100)
	push := NewPush(pool,  3,3, consumptionResults)
	push.Start()
	for i:=0;i<1000 ;i++  {
		push.Push(send)
	}
	push.Stop()
	fmt.Println("stoped===============")
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
func consumptionResults(results []Result)  {
	fmt.Println(len(results))
}