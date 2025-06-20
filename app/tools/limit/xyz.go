package limit

import (
	"crypto/md5"
	"encoding/hex"
	"govote/app/tools/log"
	"time"

	"govote/app/model/redis_cache"

	"github.com/gin-gonic/gin"
)

const (
	MaxRequests    = 5 // 最大请求次数
	BanDuration    = 3 // 限流持续时间(秒)
	WindowDuration = 5 // 滑动窗口时间(秒)
)

func CheckXYZ(context *gin.Context) bool {
	// 获取ip和ua
	ip := context.ClientIP()
	ua := context.GetHeader("user-agent")
	log.L.Infof("ip:%s, ua:%s", ip, ua)

	// 转化为md5
	hash := md5.New()
	hash.Write([]byte(ip + ua))
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	// 检查是否被限流
	banKey := "ban-" + hashString
	banFlag, err := redis_cache.GetRedisClient().Get(context, banKey).Bool()
	if err == nil && banFlag {
		log.L.Infof("IP %s 正在被限流", ip)
		return false
	}

	// 使用原子操作检查和递增计数器
	counterKey := "xyz-" + hashString
	pipe := redis_cache.GetRedisClient().Pipeline()

	// 递增计数
	incrCmd := pipe.Incr(context, counterKey)
	// 设置过期时间（只在第一次设置）
	pipe.Expire(context, counterKey, WindowDuration*time.Second)

	_, err = pipe.Exec(context)
	if err != nil {
		log.L.Errorf("Redis操作失败: %v", err)
		// Redis错误时，为了安全起见，拒绝请求
		return false
	}

	// 获取递增后的值
	currentCount := incrCmd.Val()

	// 检查是否超过限制
	if currentCount > MaxRequests {
		// 设置限流标记
		_, err = redis_cache.GetRedisClient().SetEx(context, banKey, true, BanDuration*time.Second).Result()
		if err != nil {
			log.L.Errorf("设置限流标记失败: %v", err)
		}
		log.L.Infof("IP %s 超过限制，已限流 %d 秒", ip, BanDuration)
		return false
	}

	log.L.Infof("IP %s 请求通过，当前计数: %d", ip, currentCount)
	return true
}
