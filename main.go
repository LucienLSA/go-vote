package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Print("香香编程喵喵喵！")
	g := gin.Default()
	g.GET("/login", func(context *gin.Context) {
		fmt.Print("香香编程喵喵喵！")
	})
	if err := g.Run(":8080"); err != nil {
		fmt.Print("启动失败！")
	}

}
