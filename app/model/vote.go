package model

import (
	"fmt"
	"govote/app/types"
)

func CreateVote(uid int64, vid int64, void int64) error {
	voteRecord := VoteOptUser{
		UserId:    uid,
		VoteId:    vid,
		VoteOptId: void,
	}
	err := DB.Create(&voteRecord).Error
	if err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return nil
}

func GetVotes() []Vote {
	var ret []Vote
	err := DB.Table("vote").Find(&ret).Error
	if err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return ret
}

func GetVote(voteInfo *types.Vote) Vote {
	var ret Vote
	err := DB.Table("vote").Where("id = ?", voteInfo.Id).First(&ret).Error
	if err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return ret
}

func GetVoteOpts(voteId int64) []VoteOpt {
	var voteOpts []VoteOpt
	err := DB.Where("vote_id = ?", voteId).Find(&voteOpts).Error
	if err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return voteOpts
}

func CheckUserVote(id int64, vid int64, void int64) VoteOptUser {
	var existingRecord VoteOptUser
	err := DB.Where("user_id = ? AND vote_id = ? AND vote_opt_id = ?",
		id, vid, void).First(&existingRecord).Error
	if err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return existingRecord
}
