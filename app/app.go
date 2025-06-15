package app

import (
	"govote/app/model"
	"govote/app/router"
	"govote/app/tools/logger"
)

// Start 启动器方法
func Start() {
	logger.NewLogger()
	model.NewMysql()
	defer func() {
		model.Close()
	}()
	// schedule.Start()
	router.New()
}
