package test

import (
	"govote/app/model"
	"testing"
)

func TestGetUserV1(t *testing.T) {
	// 初始化数据库
	model.NewMysql()
	defer model.Close()

	// 测试用例
	testCases := []struct {
		name        string
		username    string
		expectFound bool
	}{
		{"存在的用户", "admin", true},
		{"存在的用户2", "user1", true},
		{"不存在的用户", "nonexistent", false},
		{"空用户名", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user, err := model.GetUserV1(tc.username)

			if tc.expectFound {
				// 期望找到用户
				if err != nil {
					t.Errorf("期望找到用户 %s，但发生错误: %s", tc.username, err)
					return
				}
				if user == nil {
					t.Errorf("期望找到用户 %s，但返回nil", tc.username)
					return
				}
				if user.Id == 0 {
					t.Errorf("期望找到用户 %s，但用户ID为0", tc.username)
					return
				}
				if user.Name != tc.username {
					t.Errorf("期望用户名 %s，但得到 %s", tc.username, user.Name)
				}
				t.Logf("成功找到用户: %+v", user)
			} else {
				// 期望找不到用户
				if err == nil && user != nil && user.Id > 0 {
					t.Errorf("期望找不到用户 %s，但找到了用户: %+v", tc.username, user)
				}
				t.Logf("正确未找到用户: %s", tc.username)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	// 初始化数据库
	model.NewMysql()
	defer model.Close()

	// 测试用例
	testCases := []struct {
		name        string
		username    string
		expectFound bool
	}{
		{"存在的用户", "admin", true},
		{"存在的用户2", "user1", true},
		{"不存在的用户", "nonexistent", false},
		{"空用户名", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user, err := model.GetUser(tc.username)

			if tc.expectFound {
				// 期望找到用户
				if err != nil {
					t.Errorf("期望找到用户 %s，但发生错误: %s", tc.username, err)
					return
				}
				if user.Id == 0 {
					t.Errorf("期望找到用户 %s，但用户ID为0", tc.username)
					return
				}
				if user.Name != tc.username {
					t.Errorf("期望用户名 %s，但得到 %s", tc.username, user.Name)
				}
				t.Logf("成功找到用户: %+v", user)
			} else {
				// 期望找不到用户
				if err == nil && user.Id > 0 {
					t.Errorf("期望找不到用户 %s，但找到了用户: %+v", tc.username, user)
				}
				t.Logf("正确未找到用户: %s", tc.username)
			}
		})
	}
}

func TestGetUserPerformance(t *testing.T) {
	// 初始化数据库
	model.NewMysql()
	defer model.Close()

	// 性能测试：多次查询同一用户
	username := "admin"

	for i := 0; i < 10; i++ {
		user, err := model.GetUser(username)
		if err != nil {
			t.Errorf("第%d次查询失败: %s", i+1, err)
			continue
		}
		if user.Id == 0 {
			t.Errorf("第%d次查询未找到用户", i+1)
		}
	}

	t.Log("性能测试完成：10次查询同一用户")
}
