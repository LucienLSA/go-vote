package insert

import (
	"govote/app/model"
	"testing"
	"time"
)

func TestInsertVoteData(t *testing.T) {
	// 初始化数据库
	err := model.Init()
	if err != nil {
		t.Fatalf("数据库初始化失败: %s", err)
	}
	defer model.Close()

	// 确保所有表都已创建
	err = model.DB.AutoMigrate(&model.Vote{}, &model.VoteOpt{}, &model.VoteOptUser{})
	if err != nil {
		t.Fatalf("表迁移失败: %s", err)
	}

	t.Log("开始插入测试数据...")

	// 1. 插入投票数据
	testVotes := []model.Vote{
		{
			Title:       "你最喜欢的编程语言是什么？",
			Type:        0,     // 单选
			Status:      0,     // 正常
			Time:        86400, // 1天有效期
			UserId:      1,     // admin用户
			CreatedTime: time.Now(),
			UpdatedTime: time.Now(),
		},
		{
			Title:       "你最喜欢的Web框架是什么？",
			Type:        0,      // 单选
			Status:      0,      // 正常
			Time:        172800, // 2天有效期
			UserId:      1,      // admin用户
			CreatedTime: time.Now(),
			UpdatedTime: time.Now(),
		},
		{
			Title:       "你最喜欢的数据库是什么？",
			Type:        1,      // 多选
			Status:      0,      // 正常
			Time:        259200, // 3天有效期
			UserId:      2,      // test用户
			CreatedTime: time.Now(),
			UpdatedTime: time.Now(),
		},
		{
			Title:       "你最喜欢的操作系统是什么？",
			Type:        0,                               // 单选
			Status:      1,                               // 超时
			Time:        3600,                            // 1小时有效期（已过期）
			UserId:      3,                               // user1用户
			CreatedTime: time.Now().Add(-24 * time.Hour), // 24小时前创建
			UpdatedTime: time.Now().Add(-24 * time.Hour),
		},
	}

	// 插入投票数据
	for i, vote := range testVotes {
		result := model.DB.Table("vote").Create(&vote)
		if result.Error != nil {
			t.Errorf("插入投票 %d 失败: %s", i+1, result.Error)
		} else {
			t.Logf("成功插入投票: %s (ID: %d)", vote.Title, vote.Id)
		}
	}

	// 2. 插入投票选项数据
	testVoteOpts := []model.VoteOpt{
		// 编程语言投票选项
		{Name: "Go", VoteId: 1, Count: 15, CreatedTime: time.Now(), UpdatedTime: time.Now()},
		{Name: "Python", VoteId: 1, Count: 25, CreatedTime: time.Now(), UpdatedTime: time.Now()},
		{Name: "JavaScript", VoteId: 1, Count: 20, CreatedTime: time.Now(), UpdatedTime: time.Now()},
		{Name: "Java", VoteId: 1, Count: 18, CreatedTime: time.Now(), UpdatedTime: time.Now()},
		{Name: "C++", VoteId: 1, Count: 12, CreatedTime: time.Now(), UpdatedTime: time.Now()},

		// Web框架投票选项
		{Name: "Gin", VoteId: 2, Count: 30, CreatedTime: time.Now(), UpdatedTime: time.Now()},
		{Name: "Django", VoteId: 2, Count: 22, CreatedTime: time.Now(), UpdatedTime: time.Now()},
		{Name: "Express.js", VoteId: 2, Count: 28, CreatedTime: time.Now(), UpdatedTime: time.Now()},
		{Name: "Spring Boot", VoteId: 2, Count: 25, CreatedTime: time.Now(), UpdatedTime: time.Now()},
		{Name: "Laravel", VoteId: 2, Count: 15, CreatedTime: time.Now(), UpdatedTime: time.Now()},

		// 数据库投票选项
		{Name: "MySQL", VoteId: 3, Count: 35, CreatedTime: time.Now(), UpdatedTime: time.Now()},
		{Name: "PostgreSQL", VoteId: 3, Count: 28, CreatedTime: time.Now(), UpdatedTime: time.Now()},
		{Name: "MongoDB", VoteId: 3, Count: 20, CreatedTime: time.Now(), UpdatedTime: time.Now()},
		{Name: "Redis", VoteId: 3, Count: 32, CreatedTime: time.Now(), UpdatedTime: time.Now()},
		{Name: "SQLite", VoteId: 3, Count: 18, CreatedTime: time.Now(), UpdatedTime: time.Now()},

		// 操作系统投票选项
		{Name: "Windows", VoteId: 4, Count: 40, CreatedTime: time.Now().Add(-24 * time.Hour), UpdatedTime: time.Now().Add(-24 * time.Hour)},
		{Name: "macOS", VoteId: 4, Count: 25, CreatedTime: time.Now().Add(-24 * time.Hour), UpdatedTime: time.Now().Add(-24 * time.Hour)},
		{Name: "Linux", VoteId: 4, Count: 35, CreatedTime: time.Now().Add(-24 * time.Hour), UpdatedTime: time.Now().Add(-24 * time.Hour)},
	}

	// 插入投票选项数据
	for _, opt := range testVoteOpts {
		result := model.DB.Table("vote_opt").Create(&opt)
		if result.Error != nil {
			t.Errorf("插入投票选项 %s 失败: %s", opt.Name, result.Error)
		} else {
			t.Logf("成功插入投票选项: %s (ID: %d, 投票数: %d)", opt.Name, opt.Id, opt.Count)
		}
	}

	// 3. 插入用户投票记录数据
	testVoteOptUsers := []model.VoteOptUser{
		{UserId: 1, VoteId: 1, VoteOptId: 1, CreatedTime: time.Now()},  // admin选择了Go
		{UserId: 2, VoteId: 1, VoteOptId: 2, CreatedTime: time.Now()},  // test选择了Python
		{UserId: 3, VoteId: 1, VoteOptId: 3, CreatedTime: time.Now()},  // user1选择了JavaScript
		{UserId: 1, VoteId: 2, VoteOptId: 6, CreatedTime: time.Now()},  // admin选择了Gin
		{UserId: 2, VoteId: 2, VoteOptId: 7, CreatedTime: time.Now()},  // test选择了Django
		{UserId: 1, VoteId: 3, VoteOptId: 11, CreatedTime: time.Now()}, // admin选择了MySQL
		{UserId: 1, VoteId: 3, VoteOptId: 12, CreatedTime: time.Now()}, // admin也选择了PostgreSQL（多选）
		{UserId: 2, VoteId: 3, VoteOptId: 13, CreatedTime: time.Now()}, // test选择了MongoDB
		{UserId: 3, VoteId: 3, VoteOptId: 14, CreatedTime: time.Now()}, // user1选择了Redis
	}

	// 插入用户投票记录数据
	for i, record := range testVoteOptUsers {
		result := model.DB.Table("vote_opt_user").Create(&record)
		if result.Error != nil {
			t.Errorf("插入用户投票记录 %d 失败: %s", i+1, result.Error)
		} else {
			t.Logf("成功插入用户投票记录: 用户ID=%d, 投票ID=%d, 选项ID=%d", record.UserId, record.VoteId, record.VoteOptId)
		}
	}

	t.Log("测试数据插入完成！")
}

func TestInsertVoteDataValidation(t *testing.T) {
	// 初始化数据库
	err := model.Init()
	if err != nil {
		t.Fatalf("数据库初始化失败: %s", err)
	}
	defer model.Close()

	// 测试无效数据插入
	invalidVote := model.Vote{
		Title:       "", // 空标题
		Type:        2,  // 无效类型
		Status:      2,  // 无效状态
		Time:        -1, // 无效时间
		UserId:      0,  // 无效用户ID
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}

	result := model.DB.Table("vote").Create(&invalidVote)
	if result.Error != nil {
		t.Logf("无效数据插入被正确处理: %s", result.Error)
	} else {
		t.Logf("成功插入无效投票: %s (ID: %d)", invalidVote.Title, invalidVote.Id)
	}
}
