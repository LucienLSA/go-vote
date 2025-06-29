package logic

import (
	"fmt"
	"govote/app/db/model"
	"govote/app/db/mysql"
	"govote/app/middlewares"
	"govote/app/param"
	"govote/app/tools/auth"
	"govote/app/tools/e"
	"govote/app/tools/log"
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
		log.L.Warnf("参数错误, err:%s", err)
		context.JSON(http.StatusOK, e.ParamErr)
		return
	}
	fmt.Printf("user:%+v", user)
	// 判断输入是否为空
	if user.Name == "" || user.Password == "" || user.Password2 == "" {
		context.JSON(http.StatusOK, e.ParamErr)
		return
	}

	// 验证码校验
	if user.CaptchaId == "" || user.CaptchaCode == "" {
		context.JSON(http.StatusOK, e.CaptchaErr)
		return
	}
	// 检验验证码
	if !VerifyCaptcha(user.CaptchaId, user.CaptchaCode) {
		context.JSON(http.StatusOK, e.CaptchaErr)
		return
	}
	//校验两次密码是否相等
	if user.Password != user.Password2 {
		context.JSON(http.StatusOK, e.PasswordErr)
		return
	}
	// 检查用户是否存在
	_, exist, err := mysql.CheckUserExist(context, user.Name)
	if err != nil {
		log.L.Warnf("[mysql.CheckUserExist] 检查用户存在失败, err:%s", err)
		context.JSON(http.StatusOK, e.ServerErr)
		return
	}
	if exist {
		// 用户存在
		log.L.Warnf("用户已经存在")
		context.JSON(http.StatusOK, e.UserExistsErr)
		return
	}
	// 构建结构体，注意：应使用来自请求参数 user 的数据
	userInfo := &model.User{
		Uuid:        uid.GenSnowID(),
		Name:        user.Name,
		Password:    auth.EncryptV2(user.Password),
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}
	// 创建用户
	err = mysql.CreateUser(context, userInfo)
	if err != nil {
		log.L.Warnf("[mysql.CreateUser] 新增用户, err:%s", err)
		context.JSON(http.StatusOK, e.ServerErr)
		return
	}
	context.JSON(http.StatusOK, e.OK)
}

// 获取到经过中间件的jwt auth鉴权后登录的用户信息 ，进行下一步处理
func GetLoginUserID(c *gin.Context) (userID int64, err error) {
	uid, ok := c.Get(middlewares.CtxtUserIDKey)
	if !ok {
		err = e.ErrorUserNotLogin
		return
	}
	userID, ok = uid.(int64)
	if !ok {
		err = e.ErrorUserNotLogin
		return
	}
	return userID, nil
}
