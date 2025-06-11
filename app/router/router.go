package router

import (
	"fmt"
	"govote/app/logic"
	"html/template"

	"github.com/gin-gonic/gin"
)

func Init() error {
	// 创建Gin引擎
	r := gin.Default()

	// 添加自定义模板函数
	r.SetFuncMap(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	})

	// 加载HTML模板
	r.LoadHTMLGlob("app/view/*")

	// 公开路由（不需要登录）
	r.GET("/index", logic.Index)
	r.GET("/login", logic.GetLogin)
	r.POST("/login", logic.DoLogin)
	r.GET("/logout", logic.Logout)

	// 投票相关路由
	r.GET("/vote/:id", logic.GetVoteDetail)
	r.POST("/vote", logic.SubmitVote)
	r.GET("/vote/:id/result", logic.GetVoteResult)

	// 需要登录验证的路由组
	authGroup := r.Group("/")
	authGroup.Use(logic.CheckUser)
	{
		// 这里可以添加需要登录验证的路由
		// 例如：用户管理等
	}

	// 启动服务器
	if err := r.Run(":8080"); err != nil {
		fmt.Printf("服务器启动失败: %s\n", err)
		return err
	}
	return nil
}
