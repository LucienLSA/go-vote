package app

import (
	"govote/app/model"
	"govote/app/router"
)

// Start 启动器方法
func Start() {
	model.NewMysql()
	defer func() {
		model.Close()
	}()

	router.New()
}
