package query

import (
	"govote/app/model"
	"govote/app/types"
	"testing"
)

func TestGetVotes(t *testing.T) {
	// 初始化数据库
	err := model.Init()
	if err != nil {
		t.Fatalf("数据库初始化失败: %s", err)
	}
	defer model.Close()

	// 测试获取所有投票
	votes := model.GetVotes()
	t.Logf("总共获取到 %d 个投票", len(votes))

	for _, vote := range votes {
		t.Logf("投票ID: %d, 标题: %s, 类型: %s, 状态: %s, 创建人ID: %d",
			vote.Id,
			vote.Title,
			getVoteTypeText(vote.Type),
			getVoteStatusText(vote.Status),
			vote.UserId)
	}

	// 验证至少有一个投票
	if len(votes) == 0 {
		t.Error("没有获取到任何投票数据")
	}
}

func TestGetVote(t *testing.T) {
	// 初始化数据库
	err := model.Init()
	if err != nil {
		t.Fatalf("数据库初始化失败: %s", err)
	}
	defer model.Close()

	// 测试单个投票
	v := &types.Vote{
		Id: 1,
	}
	vote := model.GetVote(v)
	t.Logf("投票ID: %d, 标题: %s, 类型: %s, 状态: %s, 创建人ID: %d",
		vote.Id,
		vote.Title,
		getVoteTypeText(vote.Type),
		getVoteStatusText(vote.Status),
		vote.UserId)
}

func TestGetVoteOptions(t *testing.T) {
	// 初始化数据库
	err := model.Init()
	if err != nil {
		t.Fatalf("数据库初始化失败: %s", err)
	}
	defer model.Close()

	// 获取所有投票
	votes := model.GetVotes()

	// 测试获取特定投票的选项
	for _, vote := range votes {
		t.Logf("投票: %s", vote.Title)

		var options []model.VoteOpt
		err := model.DB.Where("vote_id = ?", vote.Id).Find(&options).Error
		if err != nil {
			t.Errorf("获取选项失败: %s", err)
			continue
		}

		t.Logf("  选项数量: %d", len(options))
		for _, opt := range options {
			t.Logf("    - %s (投票数: %d)", opt.Name, opt.Count)
		}

		// 验证每个投票至少有一个选项
		if len(options) == 0 {
			t.Errorf("投票 %s 没有选项", vote.Title)
		}
	}
}

func TestGetUserVoteRecords(t *testing.T) {
	// 初始化数据库
	err := model.Init()
	if err != nil {
		t.Fatalf("数据库初始化失败: %s", err)
	}
	defer model.Close()

	// 测试获取用户投票记录
	var voteRecords []model.VoteOptUser
	err = model.DB.Find(&voteRecords).Error
	if err != nil {
		t.Fatalf("获取用户投票记录失败: %s", err)
	}

	t.Logf("总共获取到 %d 条用户投票记录", len(voteRecords))
	for _, record := range voteRecords {
		t.Logf("  - 用户ID: %d, 投票ID: %d, 选项ID: %d",
			record.UserId, record.VoteId, record.VoteOptId)
	}

	// 验证投票记录的有效性
	for _, record := range voteRecords {
		if record.UserId <= 0 {
			t.Errorf("无效的用户ID: %d", record.UserId)
		}
		if record.VoteId <= 0 {
			t.Errorf("无效的投票ID: %d", record.VoteId)
		}
		if record.VoteOptId <= 0 {
			t.Errorf("无效的选项ID: %d", record.VoteOptId)
		}
	}
}

func TestVoteStatistics(t *testing.T) {
	// 初始化数据库
	err := model.Init()
	if err != nil {
		t.Fatalf("数据库初始化失败: %s", err)
	}
	defer model.Close()

	// 统计正常状态的投票数量
	var normalVoteCount int64
	model.DB.Model(&model.Vote{}).Where("status = ?", 0).Count(&normalVoteCount)
	t.Logf("正常状态投票数量: %d", normalVoteCount)

	// 统计超时状态的投票数量
	var expiredVoteCount int64
	model.DB.Model(&model.Vote{}).Where("status = ?", 1).Count(&expiredVoteCount)
	t.Logf("超时状态投票数量: %d", expiredVoteCount)

	// 统计单选投票数量
	var singleVoteCount int64
	model.DB.Model(&model.Vote{}).Where("type = ?", 0).Count(&singleVoteCount)
	t.Logf("单选投票数量: %d", singleVoteCount)

	// 统计多选投票数量
	var multiVoteCount int64
	model.DB.Model(&model.Vote{}).Where("type = ?", 1).Count(&multiVoteCount)
	t.Logf("多选投票数量: %d", multiVoteCount)

	// 验证统计数据
	totalVotes := normalVoteCount + expiredVoteCount
	if totalVotes == 0 {
		t.Error("没有找到任何投票数据")
	}

	totalTypeVotes := singleVoteCount + multiVoteCount
	if totalTypeVotes != totalVotes {
		t.Errorf("投票类型统计不匹配: 总数=%d, 单选=%d, 多选=%d", totalVotes, singleVoteCount, multiVoteCount)
	}
}

func TestTopVoteOptions(t *testing.T) {
	// 初始化数据库
	err := model.Init()
	if err != nil {
		t.Fatalf("数据库初始化失败: %s", err)
	}
	defer model.Close()

	// 测试获取最受欢迎的投票选项
	var topOptions []model.VoteOpt
	err = model.DB.Order("count desc").Limit(5).Find(&topOptions).Error
	if err != nil {
		t.Fatalf("获取最受欢迎选项失败: %s", err)
	}

	t.Log("前5个最受欢迎的选项:")
	for i, opt := range topOptions {
		t.Logf("  %d. %s (投票数: %d)", i+1, opt.Name, opt.Count)
	}

	// 验证排序是否正确
	for i := 1; i < len(topOptions); i++ {
		if topOptions[i-1].Count < topOptions[i].Count {
			t.Errorf("排序错误: %s(%d) 应该排在 %s(%d) 前面",
				topOptions[i-1].Name, topOptions[i-1].Count,
				topOptions[i].Name, topOptions[i].Count)
		}
	}
}

func TestVoteDataIntegrity(t *testing.T) {
	// 初始化数据库
	err := model.Init()
	if err != nil {
		t.Fatalf("数据库初始化失败: %s", err)
	}
	defer model.Close()

	// 测试数据完整性
	votes := model.GetVotes()

	for _, vote := range votes {
		// 验证投票基本字段
		if vote.Title == "" {
			t.Errorf("投票ID %d 的标题为空", vote.Id)
		}

		if vote.Type != 0 && vote.Type != 1 {
			t.Errorf("投票ID %d 的类型无效: %d", vote.Id, vote.Type)
		}

		if vote.Status != 0 && vote.Status != 1 {
			t.Errorf("投票ID %d 的状态无效: %d", vote.Id, vote.Status)
		}

		if vote.UserId <= 0 {
			t.Errorf("投票ID %d 的用户ID无效: %d", vote.Id, vote.UserId)
		}

		// 验证投票选项的存在性
		var optionCount int64
		model.DB.Model(&model.VoteOpt{}).Where("vote_id = ?", vote.Id).Count(&optionCount)
		if optionCount == 0 {
			t.Errorf("投票ID %d 没有选项", vote.Id)
		}
	}
}

// 获取投票类型文本
func getVoteTypeText(voteType int) string {
	if voteType == 0 {
		return "单选"
	}
	return "多选"
}

// 获取投票状态文本
func getVoteStatusText(status int) string {
	if status == 0 {
		return "正常"
	}
	return "超时"
}
