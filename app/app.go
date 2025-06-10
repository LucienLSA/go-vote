package app

import (
	"fmt"
	"govote/app/model"
	"govote/app/router"
)

func Start() {
	// 初始化数据库
	err := model.Init()
	if err != nil {
		fmt.Printf("数据库初始化失败: %s\n", err)
		return
	}
	defer func() {
		model.Close()
	}()
	err = router.Init()
}
