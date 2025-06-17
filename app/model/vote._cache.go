package model

import (
	"context"
	"encoding/json"
	"fmt"
	"govote/app/tools/log"
	"time"
)

// 从缓存中获取投票记录
func GetVoteCache(c context.Context, id int64) VoteWithOpt {
	var ret VoteWithOpt
	key := fmt.Sprintf("key_%d", id)
	log.L.Infof("key::%s\n", key)
	voteStr, err := Rdb.Get(c, key).Result()
	if err == nil || len(voteStr) > 0 {
		//存在数据
		log.L.Info("key存在, 直接返回缓存中的数据")
		_ = json.Unmarshal([]byte(voteStr), &ret)
		return ret
	}
	vote, err := GetVoteV5(id)
	if err != nil {
		log.L.Errorf("获取投票记录详情失败, err:%s\n", err.Error())
	}
	if vote.Vote.Id > 0 {
		//写入缓存
		s, _ := json.Marshal(&vote)
		err1 := Rdb.Set(c, key, s, 3600*time.Second).Err()
		if err1 != nil {
			log.L.Errorf("写入缓存失败, err:%s\n", err1.Error())
		}
		ret = *vote
	}
	return ret
}

// 删除投票记录，用于解决脏读问题，读时更新，写时删除
func CleanVote(c context.Context, voteid int64) error {
	err := Rdb.Set(c, fmt.Sprintf("key_vote_%d", voteid), "", 0).Err()
	if err != nil {
		log.L.Errorf("删除缓存失败, err:%s\n", err)
		return err
	}
	return nil
}
