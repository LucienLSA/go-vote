package model

import (
	"fmt"
	"govote/app/tools/log"
	"time"

	"gorm.io/gorm"
)

func GetVotes() []Vote {
	ret := make([]Vote, 0)
	if err := Conn.Table("vote").Find(&ret).Error; err != nil {
		log.L.Errorf("查询投票失败, err:%s\n", err)
	}
	return ret
}

func GetVote(id int64) VoteWithOpt {
	var ret Vote
	if err := Conn.Table("vote").Where("id = ?", id).First(&ret).Error; err != nil {
		log.L.Errorf("查询投票记录失败, err:%s\n", err)
	}

	opt := make([]VoteOpt, 0)
	if err := Conn.Table("vote_opt").Where("vote_id = ?", id).Find(&opt).Error; err != nil {
		log.L.Errorf("查询投票选项失败, err:%s\n", err)
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
		log.L.Errorf("查询投票记录失败, err:%s\n", err)
		tx.Rollback()
	}

	for _, value := range optIDs {
		if err := tx.Table("vote_opt").Where("id = ?", value).Update("count", gorm.Expr("count + ?", 1)).Error; err != nil {
			log.L.Errorf("更新投票选项失败, err:%s\n", err)
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
			log.L.Errorf("创建投票记录失败, err:%s\n", err)
			tx.Rollback()
		}
	}
	tx.Commit()
	return true
}

// DoVoteV1 匿名函数形式
func DoVoteV1(userId, voteId int64, optIDs []int64) bool {
	err := Conn.Transaction(func(tx *gorm.DB) error {
		var ret Vote
		if err := tx.Table("vote").Where("id = ?", voteId).First(&ret).Error; err != nil {
			log.L.Errorf("查询投票记录失败, err:%s\n", err)
			return err //只要返回了err GORM会直接回滚，不会提交
		}

		for _, value := range optIDs {
			if err := tx.Table("vote_opt").Where("id = ?", value).Update("count", gorm.Expr("count + ?", 1)).Error; err != nil {
				log.L.Errorf("更新投票选项失败, err:%s\n", err)
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
				log.L.Errorf("创建投票记录失败, err:%s\n", err)
				return err
			}
		}
		return nil //如果返回nil 则直接commit
	})
	return err == nil
}

// 添加事务检验重复投票
func DoVoteV2(userId, voteId int64, optIDs []int64) bool {
	tx := Conn.Begin()
	var ret Vote
	if err := tx.Table("vote").Where("id = ?", voteId).First(&ret).Error; err != nil {
		log.L.Errorf("查询投票记录失败, err:%s\n", err)
		tx.Rollback()
	}

	//检查是否投过票
	var oldUser VoteOptUser
	if err := tx.Table("vote_opt_user").Where("user_id = ? and vote_id = ?", userId, voteId).First(&oldUser).Error; err != nil {
		log.L.Errorf("查询用户与投票记录失败, err:%s\n", err)
		tx.Rollback()
		return false
	}
	if oldUser.Id > 0 {
		log.L.Error("用户已经投过票!")
		tx.Rollback()
		return false
	}

	for _, value := range optIDs {
		if err := tx.Table("vote_opt").Where("id = ?", value).Update("count", gorm.Expr("count + ?", 1)).Error; err != nil {
			log.L.Errorf("更新投票选项失败, err:%s\n", err)
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
			log.L.Errorf("创建投票记录失败, err:%s\n", err)
			tx.Rollback()
		}
	}
	tx.Commit()
	return true
}

func GetVoteHistory(userId, voteId int64) []VoteOptUser {
	//检查是否投过票
	ret := make([]VoteOptUser, 0)
	if err := Conn.Table("vote_opt_user").Where("user_id = ? and vote_id = ?", userId, voteId).First(&ret).Error; err != nil {
		log.L.Errorf("查询用户与投票记录失败, err:%s\n", err)
		return nil
	}
	return ret
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

// 改造为原生SQL
func DelVoteV1(id int64) bool {
	if err := Conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`delete from vote where id = ? limit 1`, id).Error; err != nil {
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

func EndVote() {
	votes := make([]Vote, 0)
	// 查询当前投票记录的状态为1的记录
	if err := Conn.Table("vote").Where("status = ?", 1).Find(&votes).Error; err != nil {
		log.L.Errorf("查询投票记录失败, err:%s\n", err)
		return
	}
	// 判断其是否过期
	now := time.Now().Unix()
	for _, vote := range votes {
		if vote.Time+vote.CreatedTime.Unix() <= now {
			// 过期将其对应的记录状态改为0
			err := Conn.Table("vote").Where("id = ?", vote.Id).Update("status", 0).Error
			if err != nil {
				log.L.Errorf("更新投票状态失败, err:%s\n", err)
				return
			}
		}
	}
}

func EndVoteV1() {
	votes := make([]Vote, 0)
	// 查询当前投票记录的状态为1的记录
	err := Conn.Raw(`select * from vote where status = ?`, 1).Scan(&votes).Error
	if err != nil {
		log.L.Errorf("查询投票记录失败, err:%s\n", err)
		return
	}
	// 判断其是否过期
	now := time.Now().Unix()
	for _, vote := range votes {
		if vote.Time+vote.CreatedTime.Unix() <= now {
			// 过期将其对应的记录状态改为9
			err := Conn.Exec(`update vote set status = 0 where id = ? limit 1`, vote.Id).Error
			if err != nil {
				log.L.Errorf("更新投票状态失败, err:%s\n", err)
				return
			}
		}
	}
}
