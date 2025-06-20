package redis_cache

import (
	"context"
	"encoding/json"
	"fmt"
	"govote/app/model"
	"govote/app/model/mysql"
	"govote/app/tools/log"
	"time"
)

// 从缓存中获取投票记录
func GetVoteCache(c context.Context, id int64) model.VoteWithOpt {
	var ret model.VoteWithOpt

	// 检查Redis连接是否可用
	if rdb == nil {
		log.L.Warn("Redis连接未初始化，直接从数据库查询")
		vote, err := mysql.GetVoteV5(id)
		if err != nil {
			log.L.Errorf("获取投票记录详情失败, err:%s\n", err.Error())
			return ret
		}
		if vote != nil {
			ret = *vote
		}
		return ret
	}

	key := fmt.Sprintf("key_%d", id)
	log.L.Infof("key::%s\n", key)
	voteStr, err := rdb.Get(c, key).Result()
	if err == nil || len(voteStr) > 0 {
		//存在数据
		log.L.Info("key存在, 直接返回缓存中的数据")
		_ = json.Unmarshal([]byte(voteStr), &ret)
		return ret
	}
	// 缓存中不存在，从数据库获取信息
	vote, err := mysql.GetVoteV5(id)
	if err != nil {
		log.L.Errorf("获取投票记录详情失败, err:%s\n", err.Error())
	}
	if vote != nil && vote.Vote.Id > 0 {
		//写入缓存
		s, _ := json.Marshal(&vote)
		err1 := rdb.Set(c, key, s, 3600*time.Second).Err()
		if err1 != nil {
			log.L.Errorf("写入缓存失败, err:%s\n", err1.Error())
		}
		ret = *vote
	}
	return ret
}

func GetVoteUserHistory(c context.Context, userId, voteId int64) ([]model.VoteOptUser, error) {
	ret := make([]model.VoteOptUser, 0)

	// 检查Redis连接是否可用
	if rdb == nil {
		log.L.Warn("Redis连接未初始化，直接从数据库查询")
		ret, err := mysql.GetVoteUser(userId, voteId)
		if err != nil {
			log.L.Errorf("获取投票用户历史失败, err:%s\n", err.Error())
			return nil, err
		}
		return ret, nil
	}

	//先查询缓存
	k := fmt.Sprintf("vote-user-%d-%d", userId, voteId)
	str, _ := rdb.Get(c, k).Result()
	fmt.Printf("str:%s\n", str)
	if len(str) > 0 {
		//将数据转化为struct
		_ = json.Unmarshal([]byte(str), &ret)
		return ret, nil
	}
	//不存在就先查数据库再封装缓存
	ret, err := mysql.GetVoteUser(userId, voteId)
	if err != nil {
		log.L.Errorf("获取投票用户历史失败, err:%s\n", err.Error())
		return nil, err
	}

	// 查到数据库
	if len(ret) > 0 {
		s, _ := json.Marshal(ret)
		err := rdb.Set(c, k, s, 3600*time.Second).Err()
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

	err := rdb.Set(c, fmt.Sprintf("key_vote_%d", voteid), "", 0).Err()
	if err != nil {
		log.L.Errorf("删除缓存失败, err:%s\n", err)
		return err
	}
	return nil
}
