package logic

import (
	"govote/app/db/mysql"
	"govote/app/db/redis_cache"
	"govote/app/param"
	"govote/app/tools/e"
	"govote/app/tools/jwt"
	"govote/app/tools/session"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetLogin(context *gin.Context) {
	context.HTML(http.StatusOK, "login.html", nil)
}

// DoLogin godoc
// @Summary      用户登录
// @Description  用户登录
// @Tags         login
// @Accept       json
// @Produce      json
// @Param        name   body      param.UserData true	"login param.UserData"
// @Success      200  {object}  e.ECode
// @Router       /login [post]
func DoLogin(context *gin.Context) {
	var user param.UserData
	if err := context.ShouldBind(&user); err != nil {
		context.JSON(http.StatusOK, e.ParamErr)
		return
	}

	// 检查用户名和密码是否为空
	if user.Name == "" || user.Password == "" {
		context.JSON(http.StatusOK, e.UserErr)
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

	ret, err := mysql.GetUser(context, user.Name)
	if err != nil {
		// 用户不存在或获取失败
		context.JSON(http.StatusOK, e.UserErr)
		return
	}

	// 使用 bcrypt.CompareHashAndPassword 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(ret.Password), []byte(user.Password)); err != nil {
		// 密码不匹配
		context.JSON(http.StatusOK, e.UserErr)
		return
	}
	// 使用sesstion保存登录的用户和id信息
	// session.SetSessionV1(context, user.Name, ret.Id)
	// 使用jwt
	token, err := jwt.GenToken(ret.Id, user.Name)
	if err != nil {
		context.JSON(http.StatusOK, e.ParamErr)
		return
	}
	ret.Token = token
	// 将token保存到redis中
	err = redis_cache.StorgeUserIdToken(token, ret.Name)
	if err != nil {
		context.JSON(http.StatusOK, e.ParamErr)
		return
	}
	context.JSON(http.StatusOK, e.ECode{
		Code:    0,
		Message: "登录成功",
		Data:    token,
	})
}

// Logout godoc
// @Summary      用户登出
// @Description  用户登出
// @Tags         login
// @Accept       json
// @Produce      json
// @Success      200  {object}  e.ECode
// @Router       /logout [post]
func Logout(context *gin.Context) {
	session.FlushSessionV1(context)
	// 检查请求的Content-Type，如果是JSON请求则返回JSON响应
	if context.GetHeader("Content-Type") == "application/json" ||
		context.Request.Method == "POST" {
		context.JSON(http.StatusOK, e.OK)
	} else {
		// 否则重定向到登录页面
		context.Redirect(http.StatusFound, "/login")
	}
}
