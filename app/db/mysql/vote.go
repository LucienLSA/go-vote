package mysql

import (
	"context"
	"govote/app/db/model"
	"govote/app/tools/log"
	"sync"
	"time"

	"gorm.io/gorm"
)

func GetVotes(ctx context.Context) []model.Vote {
	ret := make([]model.Vote, 0)
	db := NewDBClient(ctx)
	if err := db.Table("vote").Find(&ret).Error; err != nil {
		log.L.Errorf("查询投票失败, err:%s\n", err)
	}
	return ret
}

// func GetVotesV1() []model.Vote {
// 	ret := make([]model.Vote, 0)
// 	ctx := context.Background()
// 	db := NewDBClient(ctx)
// 	if err := db.Raw(`select * from vote`).Scan(&ret).Error; err != nil {
// 		log.L.Errorf("查询投票失败, err:%s\n", err)
// 	}
// 	return ret
// }

func GetVote(ctx context.Context, id int64) model.VoteWithOpt {
	var ret model.Vote
	db := NewDBClient(ctx)
	if err := db.Table("vote").Where("id = ?", id).First(&ret).Error; err != nil {
		log.L.Errorf("查询投票记录失败, err:%s\n", err)
	}

	opt := make([]model.VoteOpt, 0)
	if err := db.Table("vote_opt").Where("vote_id = ?", id).Find(&opt).Error; err != nil {
		log.L.Errorf("查询投票选项失败, err:%s\n", err)
	}
	return model.VoteWithOpt{
		Vote: ret,
		Opt:  opt,
	}
}

// func GetVoteV1(id int64) model.VoteWithOpt {
// 	var ret model.Vote
// 	ctx := context.Background()
// 	db := NewDBClient(ctx)
// 	if err := db.Raw(`select * from vote where id = ?`, id).Scan(&ret).Error; err != nil {
// 		log.L.Errorf("查询投票记录失败, err:%s\n", err)
// 	}
// 	opt := make([]model.VoteOpt, 0)
// 	if err := db.Raw(`select * from vote_opt where vote_id = ?`, id).Scan(&opt).Error; err != nil {
// 		log.L.Errorf("查询投票选项失败, err:%s\n", err)
// 	}
// 	return model.VoteWithOpt{
// 		Vote: ret,
// 		Opt:  opt,
// 	}
// }

// // 预加载模式
// func GetVoteV2(id int64) (*model.Vote, error) {
// 	var ret model.Vote
// 	ctx := context.Background()
// 	db := NewDBClient(ctx)
// 	err := db.Preload("Opt").Table("vote").Where("id = ?", id).First(&ret).Error
// 	if err != nil {
// 		log.L.Errorf("查询投票记录失败, err:%s\n", err)
// 		return &ret, err
// 	}
// 	return &ret, nil
// }

// // JOIN
// func GetVoteV3(id int64) (*model.VoteWithOpt, error) {
// 	var ret model.VoteWithOpt
// 	ctx := context.Background()
// 	db := NewDBClient(ctx)
// 	sql := `select vote.*,vote_opt.id as vid, vote_opt.name,vote_opt.count from vote join vote_opt on vote.id = vote_opt.vote_id where vote.id = ?`
// 	row, err := db.Raw(sql, id).Rows()
// 	if err != nil {
// 		log.L.Errorf("查询投票记录失败, err:%s\n", err)
// 		return nil, err
// 	}
// 	for row.Next() {
// 		temp := make(map[string]any)
// 		_ = db.ScanRows(row, &temp)
// 		fmt.Printf("temp:%+v\n", temp)
// 		if v, ok := temp["id"]; ok {
// 			ret.Vote.Id = v.(int64)
// 		}
// 	}
// 	return &ret, nil
// }

// // 协程 并发
// func GetVoteV4(id int64) (*model.VoteWithOpt, error) {
// 	var ret model.Vote
// 	ctx := context.Background()
// 	db := NewDBClient(ctx)
// 	ch := make(chan struct{}, 2)
// 	go func() {
// 		if err := db.Table("vote").Where("id = ?", id).First(&ret).Error; err != nil {
// 			log.L.Errorf("查询投票记录失败, err:%s\n", err)
// 		}
// 		ch <- struct{}{}
// 	}()
// 	opt := make([]model.VoteOpt, 0)
// 	go func() {
// 		if err := db.Table("vote_opt").Where("vote_id = ?", id).Find(&opt).Error; err != nil {
// 			log.L.Errorf("查询投票选项失败, err:%s\n", err)
// 		}
// 		ch <- struct{}{}
// 	}()
// 	var cnt int
// 	for _ = range ch {
// 		cnt++
// 		if cnt >= 2 {
// 			break
// 		}
// 	}
// 	return &model.VoteWithOpt{
// 		Vote: ret,
// 		Opt:  opt,
// 	}, nil
// }

// waitGroup
func GetVoteV5(ctx context.Context, id int64) (*model.VoteWithOpt, error) {
	var ret model.Vote
	db := NewDBClient(ctx)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := db.Table("vote").Where("id = ?", id).First(&ret).Error; err != nil {
			log.L.Errorf("查询投票记录失败, err:%s\n", err)
		}
	}()
	opt := make([]model.VoteOpt, 0)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := db.Table("vote_opt").Where("vote_id = ?", id).Find(&opt).Error; err != nil {
			log.L.Errorf("查询投票选项失败, err:%s\n", err)
		}
	}()
	wg.Wait()
	return &model.VoteWithOpt{
		Vote: ret,
		Opt:  opt,
	}, nil
}

func GetVoteByName(ctx context.Context, name string) *model.Vote {
	var ret model.Vote
	db := NewDBClient(ctx)
	if err := db.Raw(`select * from vote where title = ?`, name).Scan(&ret).Error; err != nil {
		log.L.Errorf("查询投票记录失败, err:%s\n", err)
	}
	return &ret
}

// func DoVote(userId, voteId int64, optIDs []int64) bool {
// 	ctx := context.Background()
// 	db := NewDBClient(ctx)
// 	tx := db.Begin()
// 	defer func() {
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 		}
// 	}()
// 	var ret model.Vote
// 	if err := tx.Table("vote").Where("id = ?", voteId).First(&ret).Error; err != nil {
// 		log.L.Errorf("查询投票记录失败, err:%s\n", err)
// 		tx.Rollback()
// 		return false
// 	}
// 	for _, value := range optIDs {
// 		if err := tx.Table("vote_opt").Where("id = ?", value).Update("count", gorm.Expr("count + ?", 1)).Error; err != nil {
// 			log.L.Errorf("更新投票选项失败, err:%s\n", err)
// 			tx.Rollback()
// 			return false
// 		}
// 		user := model.VoteOptUser{
// 			VoteId:      voteId,
// 			UserId:      userId,
// 			VoteOptId:   value,
// 			CreatedTime: time.Now(),
// 		}
// 		err := tx.Create(&user).Error // 通过数据的指针来创建
// 		if err != nil {
// 			log.L.Errorf("创建投票记录失败, err:%s\n", err)
// 			tx.Rollback()
// 			return false
// 		}
// 	}
// 	if err := tx.Commit().Error; err != nil {
// 		log.L.Errorf("提交事务失败, err:%s\n", err)
// 		return false
// 	}
// 	return true
// }
// // DoVoteV1 匿名函数形式
// func DoVoteV1(userId, voteId int64, optIDs []int64) bool {
// 	ctx := context.Background()
// 	db := NewDBClient(ctx)
// 	err := db.Transaction(func(tx *gorm.DB) error {
// 		var ret model.Vote
// 		if err := tx.Table("vote").Where("id = ?", voteId).First(&ret).Error; err != nil {
// 			log.L.Errorf("查询投票记录失败, err:%s\n", err)
// 			return err //只要返回了err GORM会直接回滚，不会提交
// 		}
// 		for _, value := range optIDs {
// 			if err := tx.Table("vote_opt").Where("id = ?", value).Update("count", gorm.Expr("count + ?", 1)).Error; err != nil {
// 				log.L.Errorf("更新投票选项失败, err:%s\n", err)
// 				return err
// 			}
// 			user := model.VoteOptUser{
// 				VoteId:      voteId,
// 				UserId:      userId,
// 				VoteOptId:   value,
// 				CreatedTime: time.Now(),
// 			}
// 			err := tx.Create(&user).Error // 通过数据的指针来创建
// 			if err != nil {
// 				log.L.Errorf("创建投票记录失败, err:%s\n", err)
// 				return err
// 			}
// 		}
// 		return nil //如果返回nil 则直接commit
// 	})
// 	return err == nil
// }

// 添加事务检验重复投票
func DoVoteV2(ctx context.Context, userId, voteId int64, optIDs []int64) error {
	db := NewDBClient(ctx)
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var ret model.Vote
	if err := tx.Table("vote").Where("id = ?", voteId).First(&ret).Error; err != nil {
		log.L.Errorf("查询投票记录失败, err:%s\n", err)
		tx.Rollback()
		return err
	}
	//检查是否投过票
	var oldUser model.VoteOptUser
	if err := tx.Table("vote_opt_user").Where("user_id = ? and vote_id = ?", userId, voteId).First(&oldUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 用户没投过票，这是正常情况
		} else {
			log.L.Errorf("查询用户与投票记录失败, err:%s\n", err)
			tx.Rollback()
			return err
		}
	} else if oldUser.Id > 0 {
		log.L.Error("用户已经投过票!")
		tx.Rollback()
		return err
	}
	for _, value := range optIDs {
		if err := tx.Table("vote_opt").Where("id = ?", value).Update("count", gorm.Expr("count + ?", 1)).Error; err != nil {
			log.L.Errorf("更新投票选项失败, err:%s\n", err)
			tx.Rollback()
			return err
		}
		user := model.VoteOptUser{
			VoteId:      voteId,
			UserId:      userId,
			VoteOptId:   value,
			CreatedTime: time.Now(),
		}
		err := tx.Create(&user).Error // 通过数据的指针来创建
		if err != nil {
			log.L.Errorf("创建投票记录失败, err:%s\n", err)
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.L.Errorf("提交事务失败, err:%s\n", err)
		return err
	}
	return nil
}

// 查询投票记录 原生SQL改造
// func DoVoteV3(userId, voteId int64, optIDs []int64) bool {
// 	ctx := context.Background()
// 	db := NewDBClient(ctx)
// 	tx := db.Begin()
// 	defer func() {
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 		}
// 	}()
// 	var ret model.Vote
// 	if err := tx.Raw(`select * from vote where id = ?`, voteId).Scan(&ret).Error; err != nil {
// 		log.L.Errorf("查询投票记录失败, err:%s\n", err)
// 		tx.Rollback()
// 		return false
// 	}
// 	//检查是否投过票
// 	var oldUser model.VoteOptUser
// 	if err := tx.Raw(`select * from vote_opt_user where user_id = ? and vote_id = ? `, userId, voteId).Scan(&oldUser).Error; err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			// 用户没投过票，这是正常情况
// 		} else {
// 			log.L.Errorf("查询用户与投票记录失败, err:%s\n", err)
// 			tx.Rollback()
// 			return false
// 		}
// 	} else if oldUser.Id > 0 {
// 		log.L.Error("用户已经投过票!")
// 		tx.Rollback()
// 		return false
// 	}
// 	for _, value := range optIDs {
// 		if err := tx.Exec(`update vote_opt set count = count + 1 where id = ? limit 1`, value).Error; err != nil {
// 			log.L.Errorf("更新投票选项失败, err:%s\n", err)
// 			tx.Rollback()
// 			return false
// 		}
// 		user := model.VoteOptUser{
// 			VoteId:      voteId,
// 			UserId:      userId,
// 			VoteOptId:   value,
// 			CreatedTime: time.Now(),
// 		}
// 		// 	err := tx.Exec(`
// 		// INSERT INTO vote_opt_users
// 		// (vote_id, user_id, vote_opt_id, created_time)
// 		// VALUES (?, ?, ?, ?)`,
// 		// 		user.VoteId,
// 		// 		user.UserId,
// 		// 		user.VoteOptId,
// 		// 		user.CreatedTime,
// 		// 	).Error
// 		err := tx.Create(&user).Error // 通过数据的指针来创建
// 		if err != nil {
// 			log.L.Errorf("创建投票记录失败, err:%s\n", err)
// 			tx.Rollback()
// 			return false
// 		}
// 	}
// 	if err := tx.Commit().Error; err != nil {
// 		log.L.Errorf("提交事务失败, err:%s\n", err)
// 		return false
// 	}
// 	return true
// }

func GetVoteUser(ctx context.Context, userId, voteId int64) ([]model.VoteOptUser, error) {
	//检查是否投过票
	db := NewDBClient(ctx)
	ret := make([]model.VoteOptUser, 0)
	if err := db.Raw(`select * from vote_opt_user where user_id = ? and vote_id = ?`, userId, voteId).Scan(&ret).Error; err != nil {
		log.L.Errorf("查询用户与投票记录失败, err:%s\n", err)
		return nil, err
	}
	return ret, nil
}

func AddVote(ctx context.Context, vote model.Vote, opt []model.VoteOpt) error {
	db := NewDBClient(ctx)
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&vote).Error; err != nil {
			log.L.Errorf("新增投票记录失败, err:%s\n", err)
			return err
		}
		for _, voteOpt := range opt {
			voteOpt.VoteId = vote.Id
			if err := tx.Create(&voteOpt).Error; err != nil {
				log.L.Errorf("新增投票记录失败, err:%s\n", err)
				return err
			}
		}
		return nil
	})
	return err
}

func UpdateVote(ctx context.Context, vote model.Vote, opt []model.VoteOpt) error {
	db := NewDBClient(ctx)
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&vote).Error; err != nil {
			log.L.Errorf("更新投票记录失败, err:%s\n", err)
			return err
		}
		for _, voteOpt := range opt {
			if err := tx.Save(&voteOpt).Error; err != nil {
				log.L.Errorf("更新投票记录失败, err:%s\n", err)
				return err
			}
		}
		return nil
	})
	return err
}

func DelVote(ctx context.Context, id int64) bool {
	db := NewDBClient(ctx)
	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.Vote{}, id).Error; err != nil {
			log.L.Errorf("删除投票记录失败, err:%s\n", err)
			return err
		}

		if err := tx.Where("vote_id = ?", id).Delete(&model.VoteOpt{}).Error; err != nil {
			log.L.Errorf("删除投票记录失败, err:%s\n", err)
			return err
		}

		if err := tx.Where("vote_id = ?", id).Delete(&model.VoteOptUser{}).Error; err != nil {
			log.L.Errorf("删除投票记录失败, err:%s\n", err)
			return err
		}

		return nil
	}); err != nil {
		log.L.Errorf("删除投票事务执行失败, err:%s\n", err)
		return false
	}
	return true
}

// 改造为原生SQL
// func DelVoteV1(id int64) bool {
// 	ctx := context.Background()
// 	db := NewDBClient(ctx)
// 	if err := db.Transaction(func(tx *gorm.DB) error {
// 		if err := tx.Exec(`delete from vote where id = ? limit 1`, id).Error; err != nil {
// 			log.L.Errorf("删除投票记录失败, err:%s\n", err)
// 			return err
// 		}

// 		if err := tx.Exec(`delete from vote_opt where vote_id = ?`, id).Error; err != nil {
// 			log.L.Errorf("删除投票记录失败, err:%s\n", err)
// 			return err
// 		}

// 		if err := tx.Exec(`delete from vote_opt_user where vote_id = ?`, id).Error; err != nil {
// 			log.L.Errorf("删除投票记录失败, err:%s\n", err)
// 			return err
// 		}

// 		return nil
// 	}); err != nil {
// 		log.L.Errorf("删除投票事务执行失败, err:%s\n", err)
// 		return false
// 	}

// 	return true
// }

func EndVote(ctx context.Context) {
	db := NewDBClient(ctx)
	votes := make([]model.Vote, 0)
	// 查询当前投票记录的状态为1的记录
	if err := db.Table("vote").Where("status = ?", 1).Find(&votes).Error; err != nil {
		log.L.Errorf("查询投票记录失败, err:%s\n", err)
		return
	}
	// 判断其是否过期
	now := time.Now().Unix()
	for _, vote := range votes {
		if vote.Time+vote.CreatedTime.Unix() <= now {
			// 过期将其对应的记录状态改为0
			err := db.Table("vote").Where("id = ?", vote.Id).Update("status", 0).Error
			if err != nil {
				log.L.Errorf("更新投票状态失败, err:%s\n", err)
				return
			}
		}
	}
}

// 原生sql
// func EndVoteV1() {
// 	ctx := context.Background()
// 	db := NewDBClient(ctx)
// 	votes := make([]model.Vote, 0)
// 	// 查询当前投票记录的状态为1的记录
// 	err := db.Raw(`select * from vote where status = ?`, 1).Scan(&votes).Error
// 	if err != nil {
// 		log.L.Errorf("查询投票记录失败, err:%s\n", err)
// 		return
// 	}
// 	// 判断其是否过期
// 	now := time.Now().Unix()
// 	for _, vote := range votes {
// 		if vote.Time+vote.CreatedTime.Unix() <= now {
// 			// 过期将其对应的记录状态改为9
// 			err := db.Exec(`update vote set status = 0 where id = ? limit 1`, vote.Id).Error
// 			if err != nil {
// 				log.L.Errorf("更新投票状态失败, err:%s\n", err)
// 				return
// 			}
// 		}
// 	}
// }
