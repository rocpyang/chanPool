package chanPool

import (
	"fmt"
	"testing"
)

func TestNewPool(t *testing.T) {
	pool := NewPool(3, 3)
	pool.Start()
	pool.AddJob(job_.do)
	pool.AddJob(job_.do)
	pool.AddJob(job_.do)
	pool.AddJob(job_.do)
	pool.AddJob(job_.do)
	pool.AddJob(job_.do)
	//time.Sleep(time.Second)
	//
	////fmt.Println("   pool.AddJob(do)", pool.AddJob(do))
	////pool.Start()
	pool.Stop()
	//time.Sleep(time.Second * 5)
	//pool.Close()
	pool.Start()
	pool.AddJob(job_.do)
	pool.AddJob(job_.do)
	pool.AddJob(job_.do)
	pool.AddJob(job_.do)
	pool.AddJob(job_.do)
	pool.AddJob(job_.do)
	pool.Stop()
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

func (job job) do() {
	fmt.Println("name=========", job.name)
	//fmt.Println("val=========",job.val)
}
