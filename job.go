package chanPool

// 工作
type Job func()

type DefaultJob interface {
	Job()error
}
