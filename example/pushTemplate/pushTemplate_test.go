package pushTemplate

import (
	"testing"
	"chanPool"
	"fmt"
)

func TestPush_Push(t *testing.T) {
	send := NewSend(&message{}, "2010005")
	poll := chanPool.NewPool(50, 20)
	push := NewPush(poll, 20, resultFun)
	push.Start()
	fmt.Println("123")
	push.Push(send)
	push.Push(send)
	push.Push(send)
	push.Push(send)
	push.Push(send)
	push.Push(send)
	push.Push(send)
	push.Push(send)
	push.Push(send)
	//time.Sleep(time.Second*3)
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
