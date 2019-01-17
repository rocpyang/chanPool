package pushTemplate

type Result interface {
	IsError() bool // 是否反馈异常
	ErrMsg() string
	Message() Message
	ErrCode() string
	MsgId() int64
}

type result struct {
	message Message
	errMsg  string
	errCode string
	isErr   bool
	msgId   int64
}

func (this *result) IsError() bool {
	return this.isErr
}
func (this *result) ErrCode() string {
	return this.errCode
}
func (this *result) ErrMsg() string {
	return this.errMsg
}

func (this *result) MsgId() int64 {
	return this.msgId
}
func (this *result) Message() Message {
	return this.message
}

func errorResult(errCode, errMsg string, message Message) Result {
	return &result{
		errMsg:  errMsg,
		errCode: errCode,
		isErr:   true,
		message: message,
	}
}
func successResult(msgId int64, message Message) Result {
	return &result{
		isErr:   false,
		msgId:   msgId,
		message: message,
	}
}
