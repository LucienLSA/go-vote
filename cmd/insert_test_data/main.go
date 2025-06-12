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

	fmt.Println("测试数据插入完成！")
	fmt.Println("现在可以启动应用程序进行测试了。")
}
