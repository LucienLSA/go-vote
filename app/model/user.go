package model

import (
	"fmt"

	"gorm.io/gorm"
)

func GetUser(name string) (User, error) {
	var ret User
	err := Conn.Table("user").Where("name = ?", name).First(&ret).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ret, fmt.Errorf("用户不存在")
		}
		fmt.Printf("查询用户失败: %s", err.Error())
		return ret, err
	}
	return ret, nil
}

// CreateUser 参数是指针
func CreateUser(user *User) error {
	return Conn.Create(user).Error
}
