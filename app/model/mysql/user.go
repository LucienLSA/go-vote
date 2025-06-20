package mysql

import (
	"context"
	"govote/app/model"
	"govote/app/tools/log"

	"gorm.io/gorm"
)

// func CreateUser

func GetUser(ctx context.Context, name string) (model.User, error) {
	var ret model.User
	db := NewDBClient(ctx)
	err := db.Table("user").Where("name = ?", name).First(&ret).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.L.Errorf("gorm.ErrRecordNotFound,err:%s\n", err)
			return ret, err
		}
		log.L.Errorf("查询用户失败, err:%s\n", err)
		return ret, err
	}
	return ret, nil
}

// 原生SQL改造
func GetUserV1(ctx context.Context, name string) (*model.User, error) {
	var ret model.User
	db := NewDBClient(ctx)
	err := db.Raw(`select * from user where name = ? limit 1`, name).Scan(&ret).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.L.Errorf("gorm.ErrRecordNotFound,err:%s\n", err)
			return &ret, err
		}
		log.L.Errorf("查询用户失败, err:%s\n", err)
		return &ret, err
	}
	return &ret, nil
}

func CreateUser(ctx context.Context, user *model.User) error {
	db := NewDBClient(ctx)
	err := db.Create(user).Error
	if err != nil {
		log.L.Errorf("创建用户失败, err:%s\n", err)
		return err
	}
	return nil
}
