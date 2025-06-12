package main

import (
	"fmt"
	"govote/app/test"
)

func main() {
	fmt.Println("=== 投票系统测试数据插入工具 ===")

	// 插入测试数据
	test.InsertTestData()

	// 打印数据统计
	test.PrintTestData()

	// 打印登录测试信息
	printLoginTestInfo()

	fmt.Println("\n测试数据插入完成！")
	fmt.Println("现在可以启动应用程序进行测试了。")
	fmt.Println("\n登录测试:")
	fmt.Println("1. 访问 /login 页面")
	fmt.Println("2. 使用测试用户登录 (admin/123456, user1/123456 等)")
	fmt.Println("3. 测试各种登录场景")
}

func printLoginTestInfo() {
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
