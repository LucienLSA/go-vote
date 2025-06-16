package app

import (
	"govote/app/model"
	"govote/app/router"
	"govote/app/tools/log"
)

// Start 启动器方法
func Start() {
	log.NewLogger()
	log.L.Info("日志初始化成功!")
	model.NewMysql()
	log.L.Info("MySQL初始化成功!")
	defer func() {
		model.Close()
	}()
	// schedule.Start()
	router.New()
	log.L.Info("路由初始化成功!")
}
