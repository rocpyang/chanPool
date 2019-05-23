package chanPool

import (
	"fmt"
	"sync"
	"testing"
)

func TestNewPool(t *testing.T) {
	pool := NewPool(20, 3)
	pool.Start()
	//pool.AddJob(job_.do)
	//pool.AddJob(job_.do)
	//pool.AddJob(job_.do)
	//pool.AddJob(job_.do)
	//pool.AddJob(job_.do)
	//pool.AddJob(job_.do)
	//time.Sleep(time.Second)
	//
	////fmt.Println("   pool.AddJob(do)", pool.AddJob(do))
	////pool.Start()
	for i:=0;i<10000 ; i++ {
		//pool.AddJob(job_.do1)
		pool.AddJob(job_.do2)
	}
	pool.Stop()
	//pool.Close()
	pool.Start()
	for i:=0;i<10000 ; i++ {
		//pool.AddJob(job_.do1)
		pool.AddJob(job_.do2)
	}
	pool.Stop()
	//time.Sleep(time.Second * 5)
	//pool.Close()
	//pool.Start()
	//pool.AddJob(job_.do)
	//pool.AddJob(job_.do)
	//pool.AddJob(job_.do)
	//pool.AddJob(job_.do)
	//pool.AddJob(job_.do)
	//pool.AddJob(job_.do)
	//pool.Stop()
	//time.Sleep(time.Second * 10)
}
func do() {
	fmt.Println("=========")
}

var job_ job = job{"123", "321"}

type job struct {
	name string
	val  string
}

//func (job job) do1() {
//	resp,err:=http.Get("http://htd.xuexilicheng.com:8906/get_remote_info.php")
//	if err != nil {
//		fmt.Println("err=========", err)
//	}
//	resp.Body.Close()
//	fmt.Println("close=========")
//	//fmt.Println("val=========",job.val)
//}
var lock sync.Mutex
var num=0
func (job job) do2() {
	lock.Lock()
	defer lock.Unlock()
	num=num+1
	fmt.Println("num====",num)
	//fmt.Println("val=========",job.val)
}
