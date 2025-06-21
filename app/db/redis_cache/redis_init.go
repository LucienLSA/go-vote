package redis_cache

import (
	"context"
	"govote/app/tools/log"
	"govote/app/tools/session"

	"github.com/rbcervilla/redisstore/v9"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client
var rctx = context.Background()

func NewRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       8,  // use default DB
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
