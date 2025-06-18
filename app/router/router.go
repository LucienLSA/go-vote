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
	// 用户注册与登录
	{
		//login
		r.GET("/login", logic.GetLogin)
		r.POST("/login", logic.DoLogin)
		r.GET("/logout", logic.Logout)
		r.POST("/logout", logic.Logout)

		//register
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
		//获取投票信息
		// index.GET("/redis", func(ctx *gin.Context) {
		// 	s := model.GetVoteCache(ctx, 3)
		// 	fmt.Printf("redis:%+v\n", s)
		// })
		index.GET("/index", logic.Index) //静态页面
		index.GET("/votes", logic.GetVotes)
		index.GET("/vote", logic.GetVoteInfo)
		// 投票
		index.POST("/vote", logic.DoVote)

		// 添加、删除和更新投票
		index.POST("/vote/add", logic.AddVote)
		index.PUT("/vote/update", logic.UpdateVote)
		index.DELETE("/vote/del", logic.DelVote)

		// 获取投票结果
		index.GET("/result", logic.ResultInfo)
		index.GET("/result/info", logic.ResultVote)
	}

	if err := r.Run(":8080"); err != nil {
		log.L.Panic("gin 启动失败！")
	}
}
