package model

import (
	"errors"
	"fmt"
	"govote/app/types"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(50);not null"`
	Password string `gorm:"type:varchar(50);not null"`
}

func GetUser(userInfo *types.UserInfo) *User {
	// 查询数据库中的用户
	var user User
	err := DB.Table("user").Where("name = ?", userInfo.Name).First(&user).Error
	if err != nil {
		fmt.Printf("查询失败: %s\n", err)
	}
	if user.Password != userInfo.Password {
		fmt.Printf("查询失败: %s\n", errors.New("密码错误"))
	}
	return &user
}
