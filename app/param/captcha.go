package param

// 验证码 参数
type CaptchaData struct {
	CaptchaId string `json:"captcha_id" form:"captcha_id"`
	Data      string `json:"data" form:"data"`     // base64 图片数据
	Answer    string `json:"answer" form:"answer"` // 验证码答案
}
