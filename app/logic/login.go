package logic

import (
	"govote/app/model"
	"govote/app/tools/e"
	"govote/app/tools/session"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name        string `json:"name" form:"name"`
	Password    string `json:"password" form:"password"`
	CaptchaId   string `json:"captcha_id" form:"captcha_id"`
	CaptchaCode string `json:"captcha_code" form:"captcha_code"`
}

func GetLogin(context *gin.Context) {
	context.HTML(http.StatusOK, "login.html", nil)
}

// DoLogin godoc
// @Summary      用户登录
// @Description  用户登录
// @Tags         login
// @Accept       json
// @Produce      json
// @Param        name   body      User true	"login User"
// @Success      200  {object}  e.ECode
// @Router       /login [post]
func DoLogin(context *gin.Context) {
	var user User
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

	ret, err := model.GetUser(user.Name)
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
	session.SetSession(context, user.Name, ret.Id)
	context.JSON(http.StatusOK, e.OK)
}

func CheckUser(context *gin.Context) {
	var name string
	var id int64
	values := session.GetSession(context)
	if v, ok := values["name"]; ok {
		name = v.(string)
	}
	if v, ok := values["id"]; ok {
		id = v.(int64)
	}
	if name == "" || id < 0 {
		context.JSON(http.StatusUnauthorized, e.NotLogin)
		context.Abort()
	}
	context.Next()
}

// Logout godoc
// @Summary      用户登出
// @Description  用户登出
// @Tags         login
// @Accept       json
// @Produce      json
// @Success      200  {object}  e.ECode
// @Router       /login [get]
func Logout(context *gin.Context) {
	session.FlushSession(context)
	context.Redirect(http.StatusFound, "/login")
}
