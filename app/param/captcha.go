package param

// 投票记录 参数
type VoteData struct {
	Id int64 `json:"id" form:"id"`
}

// 投票信息 参数
type VoteInfoData struct {
	UserID int64   `json:"id" form:"id"`
	VoteID int64   `json:"vote_id" form:"vote_id"`
	Opt    []int64 `json:"opt[]" form:"opt[]"`
}

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

// 验证码 参数
type CaptchaData struct {
	CaptchaId string `json:"captcha_id" form:"captcha_id"`
	Data      string `json:"data" form:"data"`     // base64 图片数据
	Answer    string `json:"answer" form:"answer"` // 验证码答案
}
