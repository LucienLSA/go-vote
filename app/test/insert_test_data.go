package test

import (
	"fmt"
	"govote/app/model"
	"govote/app/tools/auth"
	"time"
)

// InsertTestData 插入测试数据
func InsertTestData() {
	// 初始化数据库连接
	model.NewMysql()
	defer model.Close()

	fmt.Println("开始插入测试数据...")

	// 1. 插入测试用户
	insertTestUsers()

	// 2. 插入测试投票
	insertTestVotes()

	// 3. 插入测试投票选项
	insertTestVoteOptions()

	// 4. 插入一些测试投票记录
	insertTestVoteRecords()

	fmt.Println("测试数据插入完成！")
}

// insertTestUsers 插入测试用户
func insertTestUsers() {
	fmt.Println("插入测试用户...")

	users := []struct {
		Name     string
		Password string
	}{
		{Name: "admin", Password: "123456"},
		{Name: "user1", Password: "123456"},
		{Name: "user2", Password: "123456"},
		{Name: "user3", Password: "123456"},
	}

	for _, u := range users {
		// 检查用户是否已存在，避免重复插入导致错误
		if existingUser, err := model.GetUser(u.Name); err == nil && existingUser.Id > 0 {
			fmt.Printf("用户 %s 已存在 (ID: %d)，跳过插入\n", u.Name, existingUser.Id)
			continue
		}

		newUser := model.User{
			Name:        u.Name,
			Password:    auth.EncryptV2(u.Password), // 对密码进行加密
			CreatedTime: time.Now(),
			UpdatedTime: time.Now(),
		}
		if err := model.Conn.Create(&newUser).Error; err != nil {
			fmt.Printf("插入用户失败: %s\n", err.Error())
		} else {
			fmt.Printf("成功插入用户: %s (ID: %d)\n", newUser.Name, newUser.Id)
		}
	}
}

// insertTestVotes 插入测试投票
func insertTestVotes() {
	fmt.Println("插入测试投票...")

	// 获取用户ID
	var users []model.User
	if err := model.Conn.Find(&users).Error; err != nil {
		fmt.Printf("获取用户失败: %s\n", err.Error())
		return
	}

	if len(users) == 0 {
		fmt.Println("没有用户数据，无法创建投票")
		return
	}

	votes := []model.Vote{
		{
			Title:       "你最喜欢的编程语言是什么？",
			Type:        0,     // 单选
			Status:      0,     // 正常
			Time:        86400, // 24小时
			UserId:      users[0].Id,
			CreatedTime: time.Now(),
			UpdatedTime: time.Now(),
		},
		{
			Title:       "你使用过哪些开发工具？",
			Type:        1,      // 多选
			Status:      0,      // 正常
			Time:        172800, // 48小时
			UserId:      users[0].Id,
			CreatedTime: time.Now(),
			UpdatedTime: time.Now(),
		},
		{
			Title:       "你最喜欢的操作系统是什么？",
			Type:        0,      // 单选
			Status:      0,      // 正常
			Time:        259200, // 72小时
			UserId:      users[1].Id,
			CreatedTime: time.Now(),
			UpdatedTime: time.Now(),
		},
		{
			Title:       "你希望学习哪些技术？",
			Type:        1,    // 多选
			Status:      1,    // 已结束
			Time:        3600, // 1小时
			UserId:      users[1].Id,
			CreatedTime: time.Now().Add(-24 * time.Hour), // 24小时前创建
			UpdatedTime: time.Now().Add(-23 * time.Hour), // 23小时前更新
		},
		{
			Title:       "你每天编程多长时间？",
			Type:        0,      // 单选
			Status:      0,      // 正常
			Time:        604800, // 7天
			UserId:      users[2].Id,
			CreatedTime: time.Now(),
			UpdatedTime: time.Now(),
		},
	}

	for _, vote := range votes {
		// 检查投票是否已存在，避免重复插入导致错误
		if existingVote := model.GetVote(vote.Id); existingVote.Vote.Id > 0 {
			fmt.Printf("投票 \"%s\" (ID: %d) 已存在，跳过插入\n", vote.Title, vote.Id)
			continue
		}

		if err := model.Conn.Create(&vote).Error; err != nil {
			fmt.Printf("插入投票失败: %s\n", err.Error())
		} else {
			fmt.Printf("成功插入投票: %s (ID: %d, 类型: %s, 状态: %s)\n",
				vote.Title, vote.Id,
				getVoteTypeText(vote.Type),
				getVoteStatusText(vote.Status))
		}
	}
}

// insertTestVoteOptions 插入测试投票选项
func insertTestVoteOptions() {
	fmt.Println("插入测试投票选项...")

	// 获取投票ID
	var votes []model.Vote
	if err := model.Conn.Find(&votes).Error; err != nil {
		fmt.Printf("获取投票失败: %s\n", err.Error())
		return
	}

	if len(votes) == 0 {
		fmt.Println("没有投票数据，无法创建选项")
		return
	}

	// 为每个投票创建选项
	for _, vote := range votes {
		var options []model.VoteOpt

		switch vote.Id {
		case 1: // 编程语言投票
			options = []model.VoteOpt{
				{Name: "Go", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
				{Name: "Python", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
				{Name: "JavaScript", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
				{Name: "Java", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
				{Name: "C++", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
			}
		case 2: // 开发工具投票
			options = []model.VoteOpt{
				{Name: "VS Code", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
				{Name: "IntelliJ IDEA", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
				{Name: "Vim", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
				{Name: "Sublime Text", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
				{Name: "Atom", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
			}
		case 3: // 操作系统投票
			options = []model.VoteOpt{
				{Name: "Windows", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
				{Name: "macOS", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
				{Name: "Linux", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
				{Name: "Ubuntu", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
			}
		case 4: // 学习技术投票
			options = []model.VoteOpt{
				{Name: "人工智能", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
				{Name: "区块链", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
				{Name: "云计算", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
				{Name: "大数据", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
				{Name: "物联网", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
			}
		case 5: // 编程时间投票
			options = []model.VoteOpt{
				{Name: "1-2小时", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
				{Name: "3-4小时", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
				{Name: "5-6小时", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
				{Name: "7-8小时", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
				{Name: "8小时以上", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
			}
		}

		for _, option := range options {
			// 检查选项是否已存在，避免重复插入导致错误
			var existingOpt model.VoteOpt
			if err := model.Conn.Where("name = ? AND vote_id = ?", option.Name, option.VoteId).First(&existingOpt).Error; err == nil && existingOpt.Id > 0 {
				fmt.Printf("选项 \"%s\" (投票ID: %d) 已存在，跳过插入\n", option.Name, option.VoteId)
				continue
			}
			if err := model.Conn.Create(&option).Error; err != nil {
				fmt.Printf("插入投票选项失败: %s\n", err.Error())
			} else {
				fmt.Printf("成功插入选项: %s (投票ID: %d, 选项ID: %d)\n",
					option.Name, option.VoteId, option.Id)
			}
		}
	}
}

// insertTestVoteRecords 插入一些测试投票记录
func insertTestVoteRecords() {
	fmt.Println("插入测试投票记录...")

	// 获取用户和投票选项
	var users []model.User
	var voteOpts []model.VoteOpt

	if err := model.Conn.Find(&users).Error; err != nil {
		fmt.Printf("获取用户失败: %s\n", err.Error())
		return
	}

	if err := model.Conn.Find(&voteOpts).Error; err != nil {
		fmt.Printf("获取投票选项失败: %s\n", err.Error())
		return
	}

	if len(users) == 0 || len(voteOpts) == 0 {
		fmt.Println("没有用户或投票选项数据，无法创建投票记录")
		return
	}

	// 创建一些测试投票记录
	testVotes := []struct {
		userId    int64
		voteId    int64
		voteOptId int64
	}{
		{users[1].Id, 1, voteOpts[0].Id},  // user1 投票给 Go
		{users[2].Id, 1, voteOpts[1].Id},  // user2 投票给 Python
		{users[3].Id, 1, voteOpts[2].Id},  // user3 投票给 JavaScript
		{users[1].Id, 2, voteOpts[5].Id},  // user1 选择 VS Code
		{users[1].Id, 2, voteOpts[6].Id},  // user1 选择 IntelliJ IDEA
		{users[2].Id, 2, voteOpts[5].Id},  // user2 选择 VS Code
		{users[3].Id, 3, voteOpts[8].Id},  // user3 选择 Linux
		{users[1].Id, 4, voteOpts[13].Id}, // user1 选择 人工智能
		{users[1].Id, 4, voteOpts[14].Id}, // user1 选择 区块链
		{users[2].Id, 4, voteOpts[15].Id}, // user2 选择 云计算
		{users[3].Id, 5, voteOpts[18].Id}, // user3 选择 3-4小时
	}

	for _, testVote := range testVotes {
		// 检查投票记录是否已存在，避免重复插入导致错误
		var existingRecord model.VoteOptUser
		if err := model.Conn.Where("user_id = ? AND vote_id = ? AND vote_opt_id = ?",
			testVote.userId, testVote.voteId, testVote.voteOptId).First(&existingRecord).Error; err == nil && existingRecord.Id > 0 {
			fmt.Printf("投票记录已存在: 用户%d 投票给选项%d (投票%d)，跳过插入\n",
				testVote.userId, testVote.voteOptId, testVote.voteId)
			continue
		}

		// 检查选项是否属于指定的投票
		var voteOpt model.VoteOpt
		if err := model.Conn.Where("id = ? AND vote_id = ?", testVote.voteOptId, testVote.voteId).First(&voteOpt).Error; err != nil {
			fmt.Printf("选项 %d 不属于投票 %d，跳过\n", testVote.voteOptId, testVote.voteId)
			continue
		}

		// 创建投票记录
		voteRecord := model.VoteOptUser{
			UserId:      testVote.userId,
			VoteId:      testVote.voteId,
			VoteOptId:   testVote.voteOptId,
			CreatedTime: time.Now(),
		}

		if err := model.Conn.Create(&voteRecord).Error; err != nil {
			fmt.Printf("插入投票记录失败: %s\n", err.Error())
		} else {
			fmt.Printf("成功插入投票记录: 用户%d 投票给选项%d (投票%d)\n",
				testVote.userId, testVote.voteOptId, testVote.voteId)
		}

		// 更新选项计数 (注意：这里的计数更新应该在提交投票时由业务逻辑完成，这里只是为了测试数据完整性)
		if err := model.Conn.Model(&model.VoteOpt{}).Where("id = ?", testVote.voteOptId).
			Update("count", model.Conn.Raw("count + 1")).Error; err != nil {
			fmt.Printf("更新选项计数失败: %s\n", err.Error())
		}
	}
}

// getVoteTypeText 获取投票类型文本
func getVoteTypeText(voteType int32) string {
	if voteType == 0 {
		return "单选"
	}
	return "多选"
}

// getVoteStatusText 获取投票状态文本
func getVoteStatusText(voteStatus int32) string {
	if voteStatus == 0 {
		return "正常"
	}
	return "已结束"
}

// PrintTestData 打印测试数据
func PrintTestData() {
	// 初始化数据库连接
	model.NewMysql()
	defer model.Close()

	fmt.Println("\n=== 测试数据统计 ===")

	// 统计用户
	var userCount int64
	model.Conn.Model(&model.User{}).Count(&userCount)
	fmt.Printf("用户总数: %d\n", userCount)

	// 统计投票
	var voteCount int64
	model.Conn.Model(&model.Vote{}).Count(&voteCount)
	fmt.Printf("投票总数: %d\n", voteCount)

	// 统计投票选项
	var voteOptCount int64
	model.Conn.Model(&model.VoteOpt{}).Count(&voteOptCount)
	fmt.Printf("投票选项总数: %d\n", voteOptCount)

	// 统计投票记录
	var voteRecordCount int64
	model.Conn.Model(&model.VoteOptUser{}).Count(&voteRecordCount)
	fmt.Printf("投票记录总数: %d\n", voteRecordCount)

	fmt.Println("=== 数据统计完成 ===")
}
