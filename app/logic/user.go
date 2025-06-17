package logic

import (
	"fmt"
	"govote/app/model"
	"govote/app/param"
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

// CreateUser godoc
// @Summary      用户注册
// @Description  用户注册
// @Tags         register
// @Accept       json
// @Produce      json
// @Param        name   body      param.CUserData true	"register param.CUserData"
// @Success      200  {object}  e.ECode
// @Router       /user/create [post]
func CreateUser(context *gin.Context) {
	var user param.CUserData
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
