package helpers

import (
	"govote/app/model"
	"testing"
)

// SetupTestDB 设置测试数据库
func SetupTestDB(t *testing.T) {
	err := model.Init()
	if err != nil {
		t.Fatalf("数据库初始化失败: %s", err)
	}
}

// CleanupTestDB 清理测试数据库
func CleanupTestDB(t *testing.T) {
	model.Close()
}

// SetupTestData 设置测试数据
func SetupTestData(t *testing.T) {
	// 确保所有表都已创建
	err := model.DB.AutoMigrate(&model.User{}, &model.Vote{}, &model.VoteOpt{}, &model.VoteOptUser{})
	if err != nil {
		t.Fatalf("表迁移失败: %s", err)
	}
}

// CleanupTestData 清理测试数据
func CleanupTestData(t *testing.T) {
	// 清理测试数据（按依赖关系倒序删除）
	model.DB.Exec("DELETE FROM vote_opt_user")
	model.DB.Exec("DELETE FROM vote_opt")
	model.DB.Exec("DELETE FROM vote")
	model.DB.Exec("DELETE FROM user")
}

// GetVoteTypeText 获取投票类型文本
func GetVoteTypeText(voteType int) string {
	if voteType == 0 {
		return "单选"
	}
	return "多选"
}

// GetVoteStatusText 获取投票状态文本
func GetVoteStatusText(status int) string {
	if status == 0 {
		return "正常"
	}
	return "超时"
}

// AssertVoteValid 验证投票数据有效性
func AssertVoteValid(t *testing.T, vote model.Vote) {
	if vote.Title == "" {
		t.Errorf("投票标题不能为空")
	}

	if vote.Type != 0 && vote.Type != 1 {
		t.Errorf("投票类型无效: %d", vote.Type)
	}

	if vote.Status != 0 && vote.Status != 1 {
		t.Errorf("投票状态无效: %d", vote.Status)
	}

	if vote.UserId <= 0 {
		t.Errorf("用户ID无效: %d", vote.UserId)
	}
}

// AssertVoteOptionValid 验证投票选项数据有效性
func AssertVoteOptionValid(t *testing.T, option model.VoteOpt) {
	if option.Name == "" {
		t.Errorf("选项名称不能为空")
	}

	if option.VoteId <= 0 {
		t.Errorf("投票ID无效: %d", option.VoteId)
	}

	if option.Count < 0 {
		t.Errorf("投票数量不能为负数: %d", option.Count)
	}
}

// AssertVoteRecordValid 验证投票记录数据有效性
func AssertVoteRecordValid(t *testing.T, record model.VoteOptUser) {
	if record.UserId <= 0 {
		t.Errorf("用户ID无效: %d", record.UserId)
	}

	if record.VoteId <= 0 {
		t.Errorf("投票ID无效: %d", record.VoteId)
	}

	if record.VoteOptId <= 0 {
		t.Errorf("选项ID无效: %d", record.VoteOptId)
	}
}
