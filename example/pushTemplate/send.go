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
	//req,err:=http.Get("https://www.baidu.com/")
	//if err!=nil {
	//	return
	//}
	//fmt.Println(req)
	//body,err:=ioutil.ReadAll(req.Body)
	//if err!=nil {
	//	return
	//}
	//body,err=ioutil.ReadAll(req.Body)
	//if err!=nil {
	//	return
	//}
	//fmt.Println("body=="+string(body))
	//for range body {
	//
	//}
	//for i:=0;i<10000;i++  {
	//
	//}
	return errorResult(NO_SUBSCRIBE, NO_SUBSCRIBE, this.message)
}

// 校验是否关注
func (this *send) subscribe() (result Result, ok bool) {
	//user, err := weChat.UserInfo(this.delegate_code, this.message.OpenId())
	//if err != nil {
	//	errorResult("-1", err.Error(), this.message)
	//	return
	//}
	//if user.IsError() {
	//	errorResult(user.Errcode, user.Errmsg, this.message)
	//	return
	//}
	//if user.Subscribe != 0 {
	//	ok = true
	//}
	errorResult(NO_SUBSCRIBE, NO_SUBSCRIBE, this.message)
	return
}

// 发送
func (this *send) send() (result Result) {
	//resu, err := weChat.SendTemplateMsg(this.message, this.delegate_code)
	//if err != nil {
	//	return errorResult("-2", err.Error(), this.message)
	//}
	//if resu.IsError() {
	//	return errorResult(resu.Errcode, resu.Errmsg, this.message)
	//}
	return successResult(123, this.message)
}
