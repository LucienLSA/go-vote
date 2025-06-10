package router

import (
	"fmt"
	"govote/app/logic"

	"github.com/gin-gonic/gin"
)

func Init() error {
	// 创建Gin引擎
	r := gin.Default()
	// 加载HTML模板
	r.LoadHTMLGlob("app/view/*")
	// 登录页面路由
	r.GET("/login", logic.GetLogin)
	r.POST("/login", logic.DoLogin)
	// 启动服务器
	if err := r.Run(":8080"); err != nil {
		fmt.Printf("服务器启动失败: %s\n", err)
		return err
	}
	return nil
}
