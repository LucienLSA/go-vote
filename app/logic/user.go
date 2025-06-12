package logic

import (
	"fmt"
	"govote/app/model"
	"govote/app/tools/auth"
	"govote/app/tools/e"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 常见用户的结构体
type CUser struct {
	Name      string `json:"name" form:"name"`
	Password  string `json:"password" form:"password"`
	Password2 string `json:"password_2" form:"password_2"`
}

func CreateUser(context *gin.Context) {
	var user CUser
	if err := context.ShouldBind(&user); err != nil {
		context.JSON(http.StatusOK, e.ECode{
			Code:    10001,
			Message: err.Error(), //这里有风险
		})
		return
	}
	fmt.Printf("user:%+v", user)

	//encrypt(user.Password)
	//encryptV1(user.Password)
	//encryptV2(user.Password)
	//return

	if user.Name == "" || user.Password == "" || user.Password2 == "" {
		context.JSON(http.StatusOK, e.ParamErr)
		return
	}

	//校验密码
	if user.Password != user.Password2 {
		context.JSON(http.StatusOK, e.ECode{
			Code:    10003,
			Message: "两次密码不同！", //这里有风险
		})
		return
	}

	// nameLen := len(user.Name)
	// password := len(user.Password)
	// if nameLen > 16 || nameLen < 8 || password > 16 || password < 8 {
	// 	context.JSON(http.StatusOK, e.ECode{
	// 		Code:    10005,
	// 		Message: "账号或密码大于8小于16",
	// 	})
	// 	return
	// }

	// //密码不能是纯数字 -》 数字+小写字母+大写字母
	// regex := regexp.MustCompile(`^[0-9]+$`)
	// if regex.MatchString(user.Password) {
	// 	context.JSON(http.StatusOK, e.ECode{
	// 		Code:    10006,
	// 		Message: "密码不能为纯数字", //这里有风险
	// 	})
	// 	return
	// }

	//这里有一个巨大的BUG，并发安全！
	if oldUser, _ := model.GetUser(user.Name); oldUser.Id > 0 {
		context.JSON(http.StatusOK, e.ECode{
			Code:    10004,
			Message: "用户名已存在！",
		})
		return
	}

	newUser := model.User{
		Name:        user.Name,
		Password:    auth.EncryptV2(user.Password),
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}
	if err := model.CreateUser(&newUser); err != nil {
		context.JSON(http.StatusOK, e.ECode{
			Code:    10007,
			Message: "新用户创建失败！", //这里有风险
		})
		return
	}

	context.JSON(http.StatusOK, e.OK)
}
