package logic

import (
	"fmt"

	"govote/app/model"
	"govote/app/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetLogin(context *gin.Context) {
	context.HTML(http.StatusOK, "login.html", nil)
}

func DoLogin(context *gin.Context) {
	var user types.UserInfo
	if err := context.ShouldBind(&user); err != nil {
		fmt.Printf("参数绑定失败: %s\n", err)
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "参数绑定失败",
		})
		return
	}
	data := model.GetUser(&user)
	ret := map[string]interface{}{
		"name":     data.Name,
		"password": data.Password,
	}
	// 返回查询到的用户信息
	context.JSON(http.StatusOK, ret)
}
