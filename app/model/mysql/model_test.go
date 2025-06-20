package mysql

import (
	"fmt"
	"govote/app/model"
	"testing"
	"time"
)

func TestGetVotes(t *testing.T) {
	NewMysql()
	//测试用例
	r := GetVotes()
	fmt.Printf("ret:%+v", r)
	Close()
}

func TestGetVote(t *testing.T) {
	NewMysql()
	//测试用例
	r := GetVote(1)
	fmt.Printf("ret:%+v", r)
	Close()
}

func TestGetVoteV1(t *testing.T) {
	NewMysql()
	//测试用例
	r, _ := GetVoteV5(1)
	fmt.Printf("ret:%+v", r)
	Close()
}

// func TestGetVoteV2(t *testing.T) {
// 	NewMysql()
// 	//测试用例
// 	r, _ := GetVoteV2(1)
// 	fmt.Printf("ret:%+v", r.Opt)
// 	Close()
// }

func TestDoVote(t *testing.T) {
	NewMysql()
	//测试用例
	r := DoVote(1, 1, []int64{1, 2})
	fmt.Printf("ret:%+v", r)
	Close()
}

func TestAddVote(t *testing.T) {
	NewMysql()
	vote := model.Vote{
		Title:       "测试用例",
		Type:        0,
		Status:      0,
		Time:        0,
		UserId:      0,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}
	opt := make([]model.VoteOpt, 0)
	opt = append(opt, model.VoteOpt{
		Name:        "测试选项1",
		Count:       0,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	})
	opt = append(opt, model.VoteOpt{
		Name:        "测试选项2",
		Count:       0,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	})
	r := AddVote(vote, opt)
	fmt.Printf("ret:%+v", r)
	Close()
}
