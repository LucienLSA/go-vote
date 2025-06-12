package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

func GetVotes() []Vote {
	ret := make([]Vote, 0)
	if err := Conn.Table("vote").Find(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return ret
}

func GetVote(id int64) VoteWithOpt {
	var ret Vote
	if err := Conn.Table("vote").Where("id = ?", id).First(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}

	opt := make([]VoteOpt, 0)
	if err := Conn.Table("vote_opt").Where("vote_id = ?", id).Find(&opt).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return VoteWithOpt{
		Vote: ret,
		Opt:  opt,
	}
}

func DoVote(userId, voteId int64, optIDs []int64) bool {
	tx := Conn.Begin()
	var ret Vote
	if err := tx.Table("vote").Where("id = ?", voteId).First(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
		tx.Rollback()
	}

	for _, value := range optIDs {
		if err := tx.Table("vote_opt").Where("id = ?", value).Update("count", gorm.Expr("count + ?", 1)).Error; err != nil {
			fmt.Printf("err:%s", err.Error())
			tx.Rollback()
		}
		user := VoteOptUser{
			VoteId:      voteId,
			UserId:      userId,
			VoteOptId:   value,
			CreatedTime: time.Now(),
		}
		err := tx.Create(&user).Error // 通过数据的指针来创建
		if err != nil {
			fmt.Printf("err:%s", err.Error())
			tx.Rollback()
		}
	}
	tx.Commit()
	return true
}

// DoVoteV1 匿名函数
func DoVoteV1(userId, voteId int64, optIDs []int64) bool {
	err := Conn.Transaction(func(tx *gorm.DB) error {
		var ret Vote
		if err := tx.Table("vote").Where("id = ?", voteId).First(&ret).Error; err != nil {
			fmt.Printf("err:%s", err.Error())
			return err //只要返回了err GORM会直接回滚，不会提交
		}

		for _, value := range optIDs {
			if err := tx.Table("vote_opt").Where("id = ?", value).Update("count", gorm.Expr("count + ?", 1)).Error; err != nil {
				fmt.Printf("err:%s", err.Error())
				return err
			}
			user := VoteOptUser{
				VoteId:      voteId,
				UserId:      userId,
				VoteOptId:   value,
				CreatedTime: time.Now(),
			}
			err := tx.Create(&user).Error // 通过数据的指针来创建
			if err != nil {
				fmt.Printf("err:%s", err.Error())
				return err
			}
		}
		return nil //如果返回nil 则直接commit
	})
	return err == nil
}

func AddVote(vote Vote, opt []VoteOpt) error {
	err := Conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&vote).Error; err != nil {
			return err
		}
		for _, voteOpt := range opt {
			voteOpt.VoteId = vote.Id
			if err := tx.Create(&voteOpt).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func UpdateVote(vote Vote, opt []VoteOpt) error {
	err := Conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&vote).Error; err != nil {
			return err
		}
		for _, voteOpt := range opt {
			if err := tx.Save(&voteOpt).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func DelVote(id int64) bool {
	if err := Conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&Vote{}, id).Error; err != nil {
			fmt.Printf("err:%s", err.Error())
			return err
		}

		if err := tx.Where("vote_id = ?", id).Delete(&VoteOpt{}).Error; err != nil {
			fmt.Printf("err:%s", err.Error())
			return err
		}

		if err := tx.Where("vote_id = ?", id).Delete(&VoteOptUser{}).Error; err != nil {
			fmt.Printf("err:%s", err.Error())
			return err
		}

		return nil
	}); err != nil {
		fmt.Printf("err:%s", err.Error())
		return false
	}

	return true
}
