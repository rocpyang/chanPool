package pushTemplate

const (
	PUSHTEMPLATECLOSED         = "PUSHTEMPLATECLOSED"        // 发送模板已经关闭
	NO_SUBSCRIBE               = "ERROR_NO_SUBSCRIBE"        // 未关注
	OVERDUE                    = "ERROR_OVERDUE"             // 已过期
	OFF_TIME                   = "OFF_TIME"                  // 不在发送时间段
	MORE_THAN_MAX_SEND_COUNT   = "MORE_THAN_MAX_SEND_COUNT"  // 超过最大发送次数
	MORE_THAN_MAX_SRETRY_COUNT = "MORE_THAN_MAX_RETRY_COUNT" // 超过最大发送次数
)

type Send interface {
	Send() (result Result)
}

type send struct {
	message       Message
	delegate_code string
}

func NewSend(message Message, delegate_code string) Send {
	return &send{message: message, delegate_code: delegate_code}
}

func (this *send) Send() (result Result) {

	return errorResult(NO_SUBSCRIBE, NO_SUBSCRIBE, this.message)
}

// 校验是否关注
func (this *send) subscribe() (result Result, ok bool) {

	errorResult(NO_SUBSCRIBE, NO_SUBSCRIBE, this.message)
	return
}

// 发送
func (this *send) send() (result Result) {

	return successResult(123, this.message)
}
