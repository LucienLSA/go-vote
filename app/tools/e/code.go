package e

import "fmt"

var (
	OK            = ECode{Code: 0}
	ServerErr     = ECode{Code: 10000, Message: "服务内部错误"}
	NotLogin      = ECode{Code: 10001, Message: "用户未登录"}
	ParamErr      = ECode{Code: 10002, Message: "参数错误"}
	UserErr       = ECode{Code: 10003, Message: "账号或密码错误"}
	CaptchaErr    = ECode{Code: 10004, Message: "验证码错误"}
	VoteRepeatErr = ECode{Code: 10005, Message: "投票重复"}
	PasswordErr   = ECode{Code: 10006, Message: "密码不一致"}
	UserExistsErr = ECode{Code: 10007, Message: "用户名已存在"}
	NotFoundErr   = ECode{Code: 10008, Message: "资源不存在"}
	LimitErr      = ECode{Code: 10009, Message: "请稍后重试"}
)

type ECode struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (e *ECode) String() string {
	return fmt.Sprintf("code:%d,message:%s", e.Code, e.Message)
}
