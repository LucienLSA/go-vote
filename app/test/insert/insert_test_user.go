package main

import (
	"fmt"
	"govote/app/model"
)

func main() {
	// 初始化数据库
	err := model.Init()
	if err != nil {
		fmt.Printf("数据库初始化失败: %s\n", err)
		return
	}

	// 创建测试用户
	testUsers := []model.User{
		{
			Name:     "admin",
			Password: "123456",
		},
		{
			Name:     "test",
			Password: "test123",
		},
		{
			Name:     "user1",
			Password: "password1",
		},
	}

	// 插入用户数据
	for _, user := range testUsers {
		result := model.DB.Table("user").Create(&user)
		if result.Error != nil {
			fmt.Printf("插入用户 %s 失败: %s\n", user.Name, result.Error)
		} else {
			fmt.Printf("成功插入用户: %s (ID: %d)\n", user.Name, user.ID)
		}
	}

	fmt.Println("测试用户数据插入完成！")
}
