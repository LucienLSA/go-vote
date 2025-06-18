package test

import (
	"govote/app/tools/uid"
	"strings"
	"testing"
)

func TestGetUUID(t *testing.T) {
	// 测试UUID生成
	uuid1 := uid.GetUUID()
	uuid2 := uid.GetUUID()

	// 验证UUID不为空
	if uuid1 == "" {
		t.Error("生成的UUID为空")
	}

	if uuid2 == "" {
		t.Error("生成的UUID为空")
	}

	// 验证UUID格式（标准UUID格式）
	if !isValidUUID(uuid1) {
		t.Errorf("生成的UUID格式不正确: %s", uuid1)
	}

	if !isValidUUID(uuid2) {
		t.Errorf("生成的UUID格式不正确: %s", uuid2)
	}

	// 验证两次生成的UUID不同
	if uuid1 == uuid2 {
		t.Error("两次生成的UUID相同，应该不同")
	}

	t.Logf("UUID1: %s", uuid1)
	t.Logf("UUID2: %s", uuid2)
}

func TestGenSnowID(t *testing.T) {
	// 初始化雪花算法
	if err := uid.InitSnowflake("2024-01-01", 1); err != nil {
		t.Fatalf("雪花算法初始化失败: %s", err.Error())
	}

	// 测试雪花ID生成
	id1 := uid.GenSnowID()
	id2 := uid.GenSnowID()

	// 验证ID不为0
	if id1 == 0 {
		t.Error("生成的雪花ID为0")
	}

	if id2 == 0 {
		t.Error("生成的雪花ID为0")
	}

	// 验证两次生成的ID不同
	if id1 == id2 {
		t.Error("两次生成的雪花ID相同，应该不同")
	}

	// 验证ID为正数
	if id1 < 0 {
		t.Errorf("生成的雪花ID为负数: %d", id1)
	}

	if id2 < 0 {
		t.Errorf("生成的雪花ID为负数: %d", id2)
	}

	t.Logf("雪花ID1: %d", id1)
	t.Logf("雪花ID2: %d", id2)
}

// isValidUUID 验证UUID格式是否正确
func isValidUUID(uuid string) bool {
	// UUID格式: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	parts := strings.Split(uuid, "-")
	if len(parts) != 5 {
		return false
	}

	// 检查每部分的长度
	if len(parts[0]) != 8 || len(parts[1]) != 4 || len(parts[2]) != 4 ||
		len(parts[3]) != 4 || len(parts[4]) != 12 {
		return false
	}

	// 检查是否都是十六进制字符
	for _, part := range parts {
		for _, char := range part {
			if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') ||
				(char >= 'A' && char <= 'F')) {
				return false
			}
		}
	}

	return true
}
