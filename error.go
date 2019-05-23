package chanPool

const (
	start = "START"
	stoped = "STOPED"
	closed = "CLOSED"

)

type PoolError struct {
	error string
}

func (this *PoolError) Error() string {
	return this.error
}
var START = &PoolError{error: start}
var STOPED = &PoolError{error: stoped}
var CLOSED = &PoolError{error: closed}

