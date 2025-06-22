package logic

import (
	"govote/app/db/mysql"
	"govote/app/db/redis_cache"
	"govote/app/param"
	"govote/app/tools/e"
	"govote/app/tools/jwt"
	"govote/app/tools/log"
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
		log.L.Warnf("参数错误, err:%s", err)
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
		log.L.Warnf("验证码验证失败")
		context.JSON(http.StatusOK, e.CaptchaErr)
		return
	}

	ret, err := mysql.GetUser(context, user.Name)
	if err != nil {
		// 用户不存在或获取失败
		log.L.Warnf("[mysql.GetUser] 获取用户失败, err:%s", err)
		context.JSON(http.StatusOK, e.UserErr)
		return
	}

	// 使用 bcrypt.CompareHashAndPassword 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(ret.Password), []byte(user.Password)); err != nil {
		// 密码不匹配
		log.L.Warnf("[bcrypt.CompareHashAndPassword] 验证密码失败, err:%s", err)
		context.JSON(http.StatusOK, e.UserErr)
		return
	}
	// 使用sesstion保存登录的用户和id信息
	// session.SetSessionV1(context, user.Name, ret.Id)
	// 使用jwt
	token, err := jwt.GenToken(ret.Id, user.Name)
	if err != nil {
		log.L.Warnf("[jwt.GenToken] 生成token失败, err:%s", err)
		context.JSON(http.StatusOK, e.ParamErr)
		return
	}
	ret.Token = token
	// 将token保存到redis中
	err = redis_cache.StorgeUserIdToken(token, ret.Name)
	if err != nil {
		log.L.Warnf("[redis_cache.StorgeUserIdToken] redis保存token失败, err:%s", err)
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
	// 清理 session
	session.FlushSessionV1(context)

	// 清理 Redis 里的 token
	token := context.GetHeader("Authorization")
	if token != "" {
		_ = redis_cache.DeleteUserIdToken(token) // 忽略错误
	}

	// 统一返回 JSON 响应
	context.JSON(http.StatusOK, e.ECode{
		Code:    0,
		Message: "退出登录成功",
		Data:    nil,
	})
}
