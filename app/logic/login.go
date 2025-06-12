package logic

import (
	"govote/app/model"
	"govote/app/tools/e"
	"govote/app/tools/session"
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	Name     string `json:"name" form:"name"`
	Password string `json:"password" form:"password"`
}

func GetLogin(context *gin.Context) {
	context.HTML(http.StatusOK, "login.html", nil)
}

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

	ret := model.GetUser(user.Name)
	if ret.Id < 1 || ret.Password != user.Password {
		context.JSON(http.StatusOK, e.UserErr)
		return
	}

	// 登录成功，设置Cookie
	// context.SetCookie("name", user.Name, 3600, "/", "", true, false)
	// context.SetCookie("Id", fmt.Sprint(ret.Id), 3600, "/", "", true, false)
	session.SetSession(context, user.Name, ret.Id)

	context.JSON(http.StatusOK, e.OK)
}

func CheckUser(context *gin.Context) {
	// name, err := context.Cookie("name")
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
		context.JSON(http.StatusOK, e.NotLogin)
		context.Abort()
	}
	// if err != nil || name == "" {
	// 	context.Redirect(http.StatusFound, "/login")
	// 	return
	// }
	context.Next()
}

func Logout(context *gin.Context) {
	context.SetCookie("name", "", -1, "/", "", true, false)
	context.SetCookie("Id", "", -1, "/", "", true, false)
	context.Redirect(http.StatusFound, "/login")
}
