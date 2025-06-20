package param

// 用户登录 参数
type UserData struct {
	Name        string `json:"name" form:"name"`
	Password    string `json:"password" form:"password"`
	CaptchaId   string `json:"captcha_id" form:"captcha_id"`
	CaptchaCode string `json:"captcha_code" form:"captcha_code"`
}

// 用户注册 参数
type CUserData struct {
	Name        string `json:"name" form:"name"`
	Password    string `json:"password" form:"password"`
	Password2   string `json:"password_2" form:"password_2"`
	CaptchaId   string `json:"captcha_id" form:"captcha_id"`
	CaptchaCode string `json:"captcha_code" form:"captcha_code"`
}
