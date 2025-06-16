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

// CreateUser 参数是指针
func CreateUser(user *User) error {
	err := Conn.Create(user).Error
	if err != nil {
		log.L.Errorf("创建用户失败, err:%s\n", err)
		return err
	}
	return nil
}
