package redis_cache

import (
	"govote/app/tools/e"
	"govote/app/tools/log"
	"time"

	"github.com/redis/go-redis/v9"
)

func StorgeUserIdToken(token, username string) (err error) {
	// 获取token的存活时间
	duration := time.Duration(24 * time.Hour)
	key := GetRedisKey(KeyUserIDTokenSetPrefix)
	// 存入redis
	if err = rdb.Set(rctx, key+username, token, duration).Err(); err != nil {
		log.L.Errorf("token存入redis失败, err:%s\n", err)
		return
	}
	return err
}

// 从redis取token
func GetJwtToken(username string) (token string, err error) {
	key := GetRedisKey(KeyUserIDTokenSetPrefix)
	token, err = rdb.Get(rctx, key+username).Result()
	if err == redis.Nil {
		log.L.Errorf("token为空, err:%s\n", e.ErrNotExistToken)
		return "", e.ErrNotExistToken
	}
	if err != nil {
		log.L.Errorf("token存入redis失败, err:%s\n", err)
		return
	}
	return
}
