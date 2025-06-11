package insert

import (
	"govote/app/model"
	"testing"
)

func TestInsertUser(t *testing.T) {
	// 初始化数据库
	err := model.Init()
	if err != nil {
		t.Fatalf("数据库初始化失败: %s", err)
	}
	defer model.Close()

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
			t.Errorf("插入用户 %s 失败: %s", user.Name, result.Error)
		} else {
			t.Logf("成功插入用户: %s (ID: %d)", user.Name, user.Id)
		}
	}

	t.Log("测试用户数据插入完成！")
}

func TestInsertUserWithDuplicate(t *testing.T) {
	// 初始化数据库
	err := model.Init()
	if err != nil {
		t.Fatalf("数据库初始化失败: %s", err)
	}
	defer model.Close()

	// 测试重复插入
	duplicateUser := model.User{
		Name:     "admin",
		Password: "123456",
	}

	result := model.DB.Table("user").Create(&duplicateUser)
	if result.Error == nil {
		t.Logf("成功插入重复用户: %s (ID: %d)", duplicateUser.Name, duplicateUser.Id)
	} else {
		t.Logf("重复插入被正确处理: %s", result.Error)
	}
}
