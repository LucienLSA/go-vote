package logic

import (
	"fmt"
	"govote/app/model"
	"govote/app/tools/auth"
	"govote/app/tools/e"
	"govote/app/tools/uid"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetRegister 渲染注册页面
func GetRegister(context *gin.Context) {
	context.HTML(http.StatusOK, "register.html", nil)
}

// 常见用户的结构体
type CUser struct {
	Name        string `json:"name" form:"name"`
	Password    string `json:"password" form:"password"`
	Password2   string `json:"password_2" form:"password_2"`
	CaptchaId   string `json:"captcha_id" form:"captcha_id"`
	CaptchaCode string `json:"captcha_code" form:"captcha_code"`
}

func CreateUser(context *gin.Context) {
	var user CUser
	if err := context.ShouldBind(&user); err != nil {
		context.JSON(http.StatusOK, e.ParamErr)
		return
	}
	fmt.Printf("user:%+v", user)

	if user.Name == "" || user.Password == "" || user.Password2 == "" {
		context.JSON(http.StatusOK, e.ParamErr)
		return
	}

	// 验证码校验
	if user.CaptchaId == "" || user.CaptchaCode == "" {
		context.JSON(http.StatusOK, e.CaptchaErr)
		return
	}
	if !VerifyCaptcha(user.CaptchaId, user.CaptchaCode) {
		context.JSON(http.StatusOK, e.CaptchaErr)
		return
	}

	//校验密码
	if user.Password != user.Password2 {
		context.JSON(http.StatusOK, e.PasswordErr)
		return
	}

	//这里有一个巨大的BUG，并发安全！
	if oldUser, _ := model.GetUser(user.Name); oldUser.Id > 0 {
		context.JSON(http.StatusOK, e.UserExistsErr)
		return
	}

	newUser := model.User{
		Uuid:        uid.GetUUID(),
		Name:        user.Name,
		Password:    auth.EncryptV2(user.Password),
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}
	if err := model.CreateUser(&newUser); err != nil {
		context.JSON(http.StatusOK, e.ServerErr)
		return
	}

	context.JSON(http.StatusOK, e.OK)
}
