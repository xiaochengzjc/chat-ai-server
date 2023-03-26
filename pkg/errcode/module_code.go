package errcode

var (
	ErrorAddData     = NewError(20000000, "添加数据错误")
	ErrorNoEmpty     = NewError(20000001, "不能为空")
	ErrorChatApi     = NewError(20000002, "调用chat api 错误")
	ErrorWechatLogin = NewError(20000003, "微信登录失败")
	ErrorChatBusy    = NewError(20000004, "前方拥挤，请稍后再试")
)
