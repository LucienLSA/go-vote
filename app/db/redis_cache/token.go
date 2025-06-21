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

// 删除用户token
func DeleteUserIdToken(token string) error {
	// 由于存储时是以 username 为 key，token 为 value，需要遍历所有 user key 找到对应的 token 并删除
	// 这里假设用户量不大，可以全量遍历
	pattern := GetRedisKey(KeyUserIDTokenSetPrefix) + "*"
	iter := rdb.Scan(rctx, 0, pattern, 0).Iterator()
	for iter.Next(rctx) {
		key := iter.Val()
		val, err := rdb.Get(rctx, key).Result()
		if err == nil && val == token {
			return rdb.Del(rctx, key).Err()
		}
	}
	if err := iter.Err(); err != nil {
		return err
	}
	return nil
}
