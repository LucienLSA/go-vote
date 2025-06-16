package captcha

import (
	"fmt"
	"govote/app/tools/log"

	"github.com/mojocn/base64Captcha"
)

type CaptchaData struct {
	CaptchaId string `json:"captcha_id" form:"captcha_id"`
	Data      string `json:"data" form:"data"`     // base64 图片数据
	Answer    string `json:"answer" form:"answer"` // 验证码答案
}

type driverString struct {
	Id            string
	CaptchaType   string
	VerifyValue   string
	DriverString  *base64Captcha.DriverString  //字符串
	DriverChinese *base64Captcha.DriverChinese //中文
	DriverMath    *base64Captcha.DriverMath    //数学
	DriverDigit   *base64Captcha.DriverDigit   //数字
}

// 数字驱动
var digitDriver = base64Captcha.DriverDigit{
	Height:   50,  //生成图片高度
	Width:    150, //生成图片宽度
	Length:   5,   //验证码长度
	MaxSkew:  1,   //文字的倾斜度 越大倾斜越狠，越不容易看懂
	DotCount: 1,   //背景的点数，越大，字体越模糊
}

// 使用内存驱动，相关数据会存在内存空间里
var store = base64Captcha.DefaultMemStore

func CaptchaGenerate() (CaptchaData, error) {
	var ret CaptchaData
	//注意，这里直接使用digitDriver 会报错。必须传一个指针。原因参考接口实现课程中的内容
	c := base64Captcha.NewCaptcha(&digitDriver, store)
	id, b64s, res, err := c.Generate()
	fmt.Println("验证码答案:", res) // 调试用，生产环境应该移除
	if err != nil {
		log.L.Errorf("数据表AutoMigrate失败, err:%s\n", err)
		return ret, err
	}

	ret.CaptchaId = id
	ret.Data = b64s
	ret.Answer = res
	return ret, nil
}

func CaptchaVerify(data CaptchaData) bool {
	// 直接使用 store.Verify 验证，传入验证码ID和用户输入的答案
	return store.Verify(data.CaptchaId, data.Answer, true)
}
