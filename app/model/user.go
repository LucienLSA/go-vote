package model

import (
	"govote/app/tools/log"

	"gorm.io/gorm"
)

func GetUser(name string) (User, error) {
	var ret User
	err := Conn.Table("user").Where("name = ?", name).First(&ret).Error
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
func GetUserV1(name string) (*User, error) {
	var ret User
	err := Conn.Raw(`select * from user where name = ? limit 1`, name).Scan(&ret).Error
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

func CreateUser(user *User) error {
	err := Conn.Create(user).Error
	if err != nil {
		log.L.Errorf("创建用户失败, err:%s\n", err)
		return err
	}
	return nil
}
