package logic

import (
	"fmt"

	"govote/app/model"
	"govote/app/tools/e"
	"govote/app/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetLogin(context *gin.Context) {
	// 获取错误信息（如果有的话）
	errorMsg := context.Query("error")
	data := gin.H{
		"error": errorMsg,
	}
	context.HTML(http.StatusOK, "login.html", data)
}

func DoLogin(context *gin.Context) {
	var user types.UserInfo
	if err := context.ShouldBind(&user); err != nil {
		// fmt.Printf("参数绑定失败: %s\n", err)
		// // 重定向到登录页面并传递错误信息
		// context.Redirect(http.StatusFound, "/login?error=参数绑定失败")
		context.JSON(http.StatusOK, e.ECode{
			Message: "参数绑定失败！",
		})
		return
	}

	data := model.GetUser(&user)
	if data.Id < 1 || data.Password != user.Password {
		fmt.Printf("登录失败: 用户名或密码错误\n")
		// 重定向到登录页面并传递错误信息
		// context.Redirect(http.StatusFound, "/login?error=用户名或密码错误")
		context.JSON(http.StatusOK, e.ECode{
			Message: "账号或者密码有误！",
		})
		return
	}

	// 登录成功，设置Cookie
	context.SetCookie("name", user.Name, 3600, "/", "", true, false)

	// 登录成功后重定向到首页
	// context.Redirect(http.StatusFound, "/index")
	context.JSON(http.StatusOK, e.ECode{
		Message: "登录成功！",
	})

}

func Logout(context *gin.Context) {
	// 清除Cookie
	context.SetCookie("name", "", -1, "/", "", true, false)
	// 重定向到首页
	context.Redirect(http.StatusFound, "/index")
}

func CheckUser(context *gin.Context) {
	name, err := context.Cookie("name")
	if err != nil || name == "" {
		// 登录失败重定向
		context.Redirect(http.StatusFound, "/login")
		context.Abort()
		return
	}
	context.Next()
}
