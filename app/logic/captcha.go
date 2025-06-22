package logic

import (
	"govote/app/param"
	"govote/app/tools/captcha"
	"govote/app/tools/e"
	"govote/app/tools/limit"
	"govote/app/tools/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GenerateCaptcha(context *gin.Context) {
	if !limit.CheckXYZ(context) {
		log.L.Warnf("[limit.CheckXYZ] 请求被限流")
		context.JSON(http.StatusOK, e.LimitErr)
		return
	}
	captchaData, err := captcha.CaptchaGenerate()
	if err != nil {
		log.L.Warnf("[captcha.CaptchaGenerate] 生成验证码失败")
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
	data := param.CaptchaData{
		CaptchaId: id,
		Answer:    code,
	}
	return captcha.CaptchaVerify(data)
}

// VerifyCaptchaHandler handles captcha verification
func VerifyCaptchaHandler(context *gin.Context) {
	var captchaReq param.CaptchaData
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
