package redis_cache

import (
	"context"
	"encoding/json"
	"fmt"
	"govote/app/config"
	"govote/app/db/model"
	"govote/app/db/mysql"
	"govote/app/tools/log"
	"time"

	"github.com/redis/go-redis/v9"
)

// 从缓存中获取投票记录
func GetVoteCache(c context.Context, id int64) (*model.VoteWithOpt, error) {
	if rdb == nil {
		log.L.Warn("Redis连接未初始化，直接从数据库查询")
		return mysql.GetVoteV5(c, id)
	}

	// 1. 尝试从缓存获取
	fullKey := GetRedisKey(KeyVoteSetPrefix) + fmt.Sprintf("key_%d", id)
	voteStr, err := rdb.Get(c, fullKey).Result()

	if err == nil {
		// 缓存命中，验证数据有效性
		if len(voteStr) > 0 {
			var ret model.VoteWithOpt
			if json.Unmarshal([]byte(voteStr), &ret) == nil && ret.Vote.Id > 0 {
				log.L.Infof("缓存命中，投票ID: %d", id)
				return &ret, nil // 缓存有效，直接返回
			}
		}
		// 缓存数据无效 (空或损坏)，当作未命中处理
		log.L.Warnf("缓存数据无效，投票ID: %d，将从数据库刷新", id)
	} else if err != redis.Nil {
		// 发生了真正的 Redis 错误
		log.L.Errorf("Redis GET 错误，key: %s, err: %v. 降级查询数据库", fullKey, err)
	} else {
		// 缓存未命中 (err == redis.Nil)
		log.L.Infof("缓存未命中，投票ID: %d，将查询数据库", id)
	}

	// 2. 缓存未命中或无效，查询数据库
	vote, dbErr := mysql.GetVoteV5(c, id)
	if dbErr != nil {
		log.L.Errorf("数据库查询失败, err: %s", dbErr.Error())
		return nil, dbErr
	}

	// 3. 如果数据库有数据，则回写缓存
	if vote != nil && vote.Vote.Id > 0 {
		s, marshalErr := json.Marshal(&vote)
		if marshalErr != nil {
			log.L.Errorf("JSON序列化失败，投票ID: %d, err: %v", id, marshalErr)
			return vote, nil // 即使缓存失败也返回DB数据
		}

		setErr := rdb.Set(c, fullKey, s, config.Conf.AppConfig.CacheExpireTime*time.Second).Err()
		if setErr != nil {
			log.L.Errorf("回写缓存失败, err: %s", setErr.Error())
		}
		return vote, nil
	}

	// 4. 数据库也查不到，资源确实不存在
	return nil, nil
}

func GetVoteUserHistory(c context.Context, userId, voteId int64) ([]model.VoteOptUser, error) {
	ret := make([]model.VoteOptUser, 0)
	// 检查Redis连接是否可用
	if rdb == nil {
		log.L.Warn("Redis连接未初始化，直接从数据库查询")
		ret, err := mysql.GetVoteUser(c, userId, voteId)
		if err != nil {
			log.L.Errorf("获取投票用户历史失败, err:%s\n", err.Error())
			return nil, err
		}
		return ret, nil
	}

	//先查询缓存
	k := fmt.Sprintf("vote-user-%d-%d", userId, voteId)
	// 使用前缀构建完整的Redis key
	fullKey := GetRedisKey(KeyVoteSetPrefix) + k
	str, _ := rdb.Get(c, fullKey).Result()
	fmt.Printf("str:%s\n", str)
	if len(str) > 0 {
		//将数据转化为struct
		_ = json.Unmarshal([]byte(str), &ret)
		return ret, nil
	}
	//不存在就先查数据库再封装缓存
	ret, err := mysql.GetVoteUser(c, userId, voteId)
	if err != nil {
		log.L.Errorf("获取投票用户历史失败, err:%s\n", err.Error())
		return nil, err
	}

	// 查到数据库
	if len(ret) > 0 {
		s, _ := json.Marshal(ret)
		err := rdb.Set(c, fullKey, s, config.Conf.AppConfig.CacheExpireTime*time.Second).Err()
		if err != nil {
			log.L.Errorf("设置缓存失败, err:%s\n", err.Error())
			return nil, err
		}
	}
	return ret, nil
}

// 删除投票记录，用于解决脏读问题，读时更新，写时删除
func CleanVote(c context.Context, voteid int64) error {
	if rdb == nil {
		log.L.Warn("Redis连接未初始化，跳过缓存清理")
		return nil
	}

	// 使用前缀构建完整的Redis key
	key := fmt.Sprintf("key_%d", voteid)
	fullKey := GetRedisKey(KeyVoteSetPrefix) + key
	err := rdb.Del(c, fullKey).Err()
	if err != nil {
		log.L.Errorf("删除缓存失败, err:%s\n", err)
		return err
	}
	return nil
}
