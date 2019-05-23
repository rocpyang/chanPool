package chanPool

import "errors"

const (
	start  = "START"
	stoped = "STOPED"
	closed = "CLOSED"
)

type PoolError struct {
	error      error
	defaultJob DefaultJob
}

func (this *PoolError) Error() string {
	return this.error.Error()
}
func (this *PoolError) DefaultJob() DefaultJob {
	return this.defaultJob
}

var START = &PoolError{error: errors.New(start)}
var STOPED = &PoolError{error: errors.New(stoped)}
var CLOSED = &PoolError{error: errors.New(closed)}
