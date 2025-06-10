package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	baseURL := "http://localhost:8080"

	// 测试用例
	testCases := []struct {
		name     string
		password string
		expected string
	}{
		{"admin", "123456", "登录成功"},
		{"test", "test123", "登录成功"},
		{"admin", "wrong", "密码错误"},
		{"nonexistent", "123456", "用户不存在"},
	}

	for _, tc := range testCases {
		fmt.Printf("\n=== 测试: 用户名=%s, 密码=%s ===\n", tc.name, tc.password)

		// 准备表单数据
		formData := url.Values{}
		formData.Set("name", tc.name)
		formData.Set("password", tc.password)

		// 发送POST请求
		resp, err := http.Post(baseURL+"/login", "application/x-www-form-urlencoded", strings.NewReader(formData.Encode()))
		if err != nil {
			fmt.Printf("请求失败: %s\n", err)
			continue
		}
		defer resp.Body.Close()

		// 读取响应
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("读取响应失败: %s\n", err)
			continue
		}

		fmt.Printf("状态码: %d\n", resp.StatusCode)
		fmt.Printf("响应内容: %s\n", string(body))
	}
}
