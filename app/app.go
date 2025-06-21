package app

import (
	"context"
	"govote/app/db/mysql"
	"govote/app/db/redis_cache"
	"govote/app/router"
	"govote/app/schedule"
	"govote/app/tools/log"
	"govote/app/tools/uid"
	"time"
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

	mysql.NewMysql()
	log.L.Info("MySQL初始化成功!")
	redis_cache.NewRedis()
	log.L.Info("Redis初始化成功!")
	defer func() {
		mysql.Close()
	}()
	// schedule.Start()
	router.New()
	log.L.Info("路由初始化成功!")
}

// 删除过期缓存
func StartEndVote() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	schedule.Start(ctx, 5*time.Second)
}
