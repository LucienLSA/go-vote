package redis_cache

import (
	"context"
	"fmt"
	"govote/app/config"
	"govote/app/tools/log"
	"govote/app/tools/session"

	"github.com/rbcervilla/redisstore/v9"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client
var rctx = context.Background()

func NewRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Conf.RedisConfig.Host, config.Conf.RedisConfig.Port),
		Password: config.Conf.RedisConfig.Password,
		DB:       config.Conf.RedisConfig.DB,
		PoolSize: config.Conf.RedisConfig.PoolSize,
	})
	// 初始化session
	var err error
	session.SessionStore, err = redisstore.NewRedisStore(context.TODO(), rdb)
	if err != nil {
		log.L.Panicf("初始化redisStore失败, err:%s\n", err)
	}
}

func GetRedisClient() *redis.Client {
	return rdb
}
