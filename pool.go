package chanPool

type Pool interface {
	Start()
	Stop()
	Close() error
	AddJob(job Job) error
}

type pool struct {
	dispatcher Dispatcher
}

//workerNum 工人池中的工人数量
//
//jobNum job池中的job数量
func NewPool(workerNum, jobNum int) Pool {
	return &pool{
		dispatcher: NewDispatcher(workerNum, jobNum),
	}
}

// 添加一个job到job池中
func (p *pool) AddJob(job Job) error {
	return p.dispatcher.AddJob(job)
}

// 停止所有进程
func (p *pool) Stop() {
	p.dispatcher.WaitForAll()
	p.dispatcher.Stop()
}

//  清理资源
func (p *pool) Close() error {
	return p.dispatcher.Close()
}

//Start worker pool and dispatch
func (p *pool) Start() {
	p.dispatcher.EnableWaitForAll(true)
	p.dispatcher.Start()
}
