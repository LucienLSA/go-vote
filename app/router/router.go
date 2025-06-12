package router

import (
	"govote/app/logic"

	"github.com/gin-gonic/gin"
)

func New() {
	r := gin.Default()
	r.LoadHTMLGlob("app/view/*")
	//相关的路径 放在这里
	r.GET("/", logic.Index)
	r.GET("/login", logic.GetLogin)
	r.POST("/login", logic.DoLogin)
	r.GET("/logout", logic.Logout)
	index := r.Group("")
	index.Use(logic.CheckUser)
	{
		index.GET("/index", logic.Index)
		index.GET("/vote", logic.GetVoteInfo)
		index.POST("/vote", logic.DoVote)

	}

	if err := r.Run(":8080"); err != nil {
		panic("gin 启动失败！")
	}
}
