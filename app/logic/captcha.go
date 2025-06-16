package logic

import (
	"govote/app/tools/captcha"
	"govote/app/tools/e"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GenerateCaptcha 生成验证码
func GenerateCaptcha(context *gin.Context) {
	captchaData, err := captcha.CaptchaGenerate()
	if err != nil {
		context.JSON(http.StatusOK, e.ServerErr)
		return
	}

	context.JSON(http.StatusOK, e.ECode{
		Code:    0,
		Message: "验证码生成成功",
		Data: gin.H{
			"captcha_id":    captchaData.CaptchaId,
			"captcha_image": captchaData.Data,
		},
	})
}

// VerifyCaptcha 验证验证码
func VerifyCaptcha(id, code string) bool {
	data := captcha.CaptchaData{
		CaptchaId: id,
		Answer:    code,
	}
	return captcha.CaptchaVerify(data)
}

// VerifyCaptchaHandler handles captcha verification
func VerifyCaptchaHandler(context *gin.Context) {
	var captchaReq captcha.CaptchaData
	if err := context.ShouldBind(&captchaReq); err != nil {
		context.JSON(http.StatusBadRequest, e.ParamErr)
		return
	}

	// verify the captcha
	if VerifyCaptcha(captchaReq.CaptchaId, captchaReq.Answer) {
		context.JSON(http.StatusOK, e.OK)
	} else {
		context.JSON(http.StatusOK, e.CaptchaErr)
	}
}
