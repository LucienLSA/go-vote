package redis_cache

// key健 定义常量易于查询和拆分
// 冒号分割命名空间 使用其区分不同的key
const (
	Prefix                  = "bluebell:"
	KeyUserIDTokenSetPrefix = "user:" // set; 保存登录用户及token
)

// 给redis key加上前缀
func GetRedisKey(key string) string {
	return Prefix + key
}
