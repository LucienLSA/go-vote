package e

import "errors"

var (
	ErrNotExistToken  = errors.New("不存在的key")
	ErrorUserNotLogin = errors.New("用户未登录")
)
