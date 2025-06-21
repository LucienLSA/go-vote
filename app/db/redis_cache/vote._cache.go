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
)

// 从缓存中获取投票记录
func GetVoteCache(c context.Context, id int64) (*model.VoteWithOpt, error) {
	if rdb == nil {
		log.L.Warn("Redis连接未初始化，直接从数据库查询")
		vote, err := mysql.GetVoteV5(c, id)
		if err != nil {
			log.L.Errorf("获取投票记录详情失败, err:%s\n", err.Error())
			return nil, err
		}
		if vote != nil {
			return vote, nil
		}
		return nil, nil
	}

	key1 := fmt.Sprintf("key_%d", id)
	// 使用前缀构建完整的Redis key
	fullKey := GetRedisKey(KeyVoteSetPrefix) + key1
	log.L.Infof("key::%s\n", fullKey)
	voteStr, err := rdb.Get(c, fullKey).Result()
	if err == nil || len(voteStr) > 0 {
		var ret model.VoteWithOpt
		_ = json.Unmarshal([]byte(voteStr), &ret)
		if ret.Vote.Id > 0 {
			return &ret, nil
		}
		return nil, nil
	}
	// 缓存中不存在，从数据库获取信息
	vote, err := mysql.GetVoteV5(c, id)
	if err != nil {
		log.L.Errorf("获取投票记录详情失败, err:%s\n", err.Error())
		return nil, err
	}
	if vote != nil && vote.Vote.Id > 0 {
		//写入缓存
		s, _ := json.Marshal(&vote)
		// 加入前缀key
		key := GetRedisKey(KeyVoteSetPrefix)
		err1 := rdb.Set(c, key+key1, s, config.Conf.AppConfig.CacheExpireTime*time.Second).Err()
		if err1 != nil {
			log.L.Errorf("写入缓存失败, err:%s\n", err1.Error())
		}
		return vote, nil
	}
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
	err := rdb.Set(c, fullKey, "", 0).Err()
	if err != nil {
		log.L.Errorf("删除缓存失败, err:%s\n", err)
		return err
	}
	return nil
}
