package router

import (
	"govote/app/logic"

	"github.com/gin-gonic/gin"
)

func New() {
	r := gin.Default()
	r.LoadHTMLGlob("app/view/*")
	//相关的路径 放在这里
	{
		index := r.Group("")
		index.Use(logic.CheckUser)
		//vote
		index.GET("/index", logic.Index) //静态页面

		index.GET("/votes", logic.GetVotes)
		index.GET("/vote", logic.GetVoteInfo)
		index.POST("/vote", logic.DoVote)

		index.POST("/vote/add", logic.AddVote)
		index.POST("/vote/update", logic.UpdateVote)
		index.POST("/vote/del", logic.DelVote)

		index.GET("/result", logic.ResultInfo)
		index.GET("/result/info", logic.ResultVote)
	}

	r.GET("/", logic.Index)

	{
		//login
		r.GET("/login", logic.GetLogin)
		r.POST("/login", logic.DoLogin)
		r.GET("/logout", logic.Logout)
	}

	if err := r.Run(":8080"); err != nil {
		panic("gin 启动失败！")
	}
}
