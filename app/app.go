package app

import (
	"govote/app/model"
	"govote/app/router"
	"govote/app/tools/log"
	"govote/app/tools/uid"
)

// Start 启动器方法
func Start() {
	log.NewLogger()
	log.L.Info("日志初始化成功!")

	// 初始化雪花算法
	if err := uid.InitSnowflake("2024-01-01", 1); err != nil {
		log.L.Panicf("雪花算法初始化失败, err:%s\n", err)
	}
	log.L.Info("雪花算法初始化成功!")

	model.NewMysql()
	log.L.Info("MySQL初始化成功!")
	model.NewRedis()
	log.L.Info("Redis初始化成功!")
	defer func() {
		model.Close()
	}()
	// schedule.Start()
	router.New()
	log.L.Info("路由初始化成功!")
}
