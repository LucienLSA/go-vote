package test

import (
	"fmt"
	"govote/app/model"
	"testing"
)

// TestLoginFunctionality 测试登录功能
func TestLoginFunctionality(t *testing.T) {
	// 初始化数据库
	model.NewMysql()
	defer model.Close()

	// 测试用例
	testCases := []struct {
		name     string
		username string
		password string
		expected bool
	}{
		{"正确用户名密码", "admin", "123456", true},
		{"正确用户名密码2", "user1", "123456", true},
		{"错误密码", "admin", "wrong", false},
		{"错误用户名", "nonexistent", "123456", false},
		{"空用户名", "", "123456", false},
		{"空密码", "admin", "", false},
		{"空用户名和密码", "", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user, _ := model.GetUser(tc.username)

			if tc.expected {
				// 期望登录成功
				if user.Id == 0 {
					t.Errorf("期望用户存在，但用户不存在: %s", tc.username)
				}
				if user.Password != tc.password {
					t.Errorf("期望密码匹配，但密码不匹配: %s", tc.username)
				}
			} else {
				// 期望登录失败
				if user.Id != 0 && user.Password == tc.password {
					t.Errorf("期望登录失败，但登录成功: %s", tc.username)
				}
			}
		})
	}
}

// TestUserData 测试用户数据
func TestUserData(t *testing.T) {
	// 初始化数据库
	model.NewMysql()
	defer model.Close()

	// 获取所有用户
	var users []model.User
	if err := model.Conn.Find(&users).Error; err != nil {
		t.Fatalf("获取用户数据失败: %s", err)
	}

	t.Logf("数据库中共有 %d 个用户", len(users))

	// 验证用户数据
	expectedUsers := []string{"admin", "user1", "user2", "user3"}
	userMap := make(map[string]bool)

	for _, user := range users {
		userMap[user.Name] = true
		t.Logf("用户: %s (ID: %d)", user.Name, user.Id)
	}

	// 检查是否包含所有期望的用户
	for _, expectedUser := range expectedUsers {
		if !userMap[expectedUser] {
			t.Errorf("缺少期望的用户: %s", expectedUser)
		}
	}
}

// TestLoginResponse 测试登录响应格式
func TestLoginResponse(t *testing.T) {
	// 这里可以测试登录API的响应格式
	// 由于需要HTTP服务器，这里只是示例结构

	t.Log("登录响应格式测试:")
	t.Log("- 成功登录: code=0, message='登录成功'")
	t.Log("- 失败登录: code=1, message='错误信息'")
	t.Log("- 参数错误: code=1, message='参数绑定失败！'")
	t.Log("- 空用户名: code=1, message='用户名和密码不能为空！'")
	t.Log("- 密码错误: code=1, message='账号或密码错误！'")
}

// PrintLoginTestInfo 打印登录测试信息
func PrintLoginTestInfo() {
	fmt.Println("=== 登录功能测试信息 ===")
	fmt.Println("测试用户:")
	fmt.Println("  - admin/123456")
	fmt.Println("  - user1/123456")
	fmt.Println("  - user2/123456")
	fmt.Println("  - user3/123456")
	fmt.Println()
	fmt.Println("测试场景:")
	fmt.Println("  1. 正确用户名密码登录")
	fmt.Println("  2. 错误密码登录")
	fmt.Println("  3. 不存在的用户名登录")
	fmt.Println("  4. 空用户名或密码登录")
	fmt.Println("  5. 登录成功后的Cookie设置")
	fmt.Println("  6. 登录失败的错误提示")
	fmt.Println()
	fmt.Println("响应格式:")
	fmt.Println("  成功: {\"code\":0,\"message\":\"登录成功\",\"data\":用户ID}")
	fmt.Println("  失败: {\"code\":1,\"message\":\"错误信息\"}")
	fmt.Println("=== 登录测试信息结束 ===")
}
