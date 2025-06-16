package router

import (
	"govote/app/logic"
	"govote/app/tools/log"

	_ "govote/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func New() {
	r := gin.Default()
	r.LoadHTMLGlob("app/view/*")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/", logic.Index)
	{
		//login
		r.GET("/login", logic.GetLogin)
		r.POST("/login", logic.DoLogin)
		r.GET("/logout", logic.Logout)

		//user
		r.GET("/register", logic.GetRegister)
		r.POST("/user/create", logic.CreateUser)
	}
	//验证码
	{
		r.GET("/captcha/generate", logic.GenerateCaptcha)
		r.POST("/captcha/verify", logic.VerifyCaptchaHandler)
	}
	// 改为RestFul风格接口
	index := r.Group("")
	index.Use(logic.CheckUser)
	{
		//vote
		index.GET("/index", logic.Index) //静态页面
		index.GET("/votes", logic.GetVotes)
		index.GET("/vote", logic.GetVoteInfo)

		index.POST("/vote", logic.DoVote)

		index.POST("/vote/add", logic.AddVote)
		index.PUT("/vote/update", logic.UpdateVote)
		index.DELETE("/vote/del", logic.DelVote)

		index.GET("/result", logic.ResultInfo)
		index.GET("/result/info", logic.ResultVote)
	}

	if err := r.Run(":8080"); err != nil {
		log.L.Panic("gin 启动失败！")
	}
}
