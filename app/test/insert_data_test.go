package test

// import (
// 	"fmt"
// 	"govote/app/model"
// 	"govote/app/model/mysql"
// 	"govote/app/tools/auth"
// 	"govote/app/tools/log"
// 	"govote/app/tools/uid"
// 	"testing"
// 	"time"
// )

// // TestInsertTestData 测试插入测试数据功能
// func TestInsertTestData(t *testing.T) {
// 	// 初始化日志系统
// 	log.NewLogger()

// 	// 初始化数据库连接
// 	mysql.NewMysql()
// 	defer mysql.Close()

// 	t.Log("开始插入测试数据...")

// 	// 1. 插入测试用户
// 	testInsertTestUsers(t)

// 	// 2. 插入测试投票
// 	testInsertTestVotes(t)

// 	// 3. 插入测试投票选项
// 	testInsertTestVoteOptions(t)

// 	// 4. 插入一些测试投票记录
// 	testInsertTestVoteRecords(t)

// 	t.Log("测试数据插入完成！")
// }

// // testInsertTestUsers 测试插入测试用户
// func testInsertTestUsers(t *testing.T) {
// 	// 初始化雪花算法
// 	if err := uid.InitSnowflake("2024-01-01", 1); err != nil {
// 		t.Fatalf("雪花算法初始化失败: %s", err.Error())
// 	}

// 	t.Log("开始测试用户创建...")

// 	// 测试用户数据
// 	testUsers := []struct {
// 		Name     string
// 		Password string
// 	}{
// 		{Name: "testuser1", Password: "123456"},
// 		{Name: "testuser2", Password: "123456"},
// 		{Name: "testuser3", Password: "123456"},
// 	}

// 	for _, user := range testUsers {
// 		t.Logf("尝试创建用户: %s", user.Name)

// 		// 检查用户是否已存在
// 		existingUser, err := mysql.GetUser(user.Name)
// 		if err == nil && existingUser.Id > 0 {
// 			t.Logf("用户 %s 已存在 (ID: %d)，跳过创建", user.Name, existingUser.Id)
// 			continue
// 		}
// 		// 用户不存在是正常情况，继续创建

// 		// 创建新用户
// 		newUser := model.User{
// 			Uuid:        uid.GenSnowID(),
// 			Name:        user.Name,
// 			Password:    auth.EncryptV2(user.Password),
// 			CreatedTime: time.Now(),
// 			UpdatedTime: time.Now(),
// 		}

// 		if err := mysql.CreateUser(&newUser); err != nil {
// 			t.Errorf("创建用户失败: %s", err.Error())
// 		} else {
// 			t.Logf("成功创建用户: %s (ID: %d, UUID: %d)", newUser.Name, newUser.Id, newUser.Uuid)
// 		}
// 	}

// 	t.Log("用户创建测试完成！")
// }

// // testInsertTestVotes 测试插入测试投票
// func testInsertTestVotes(t *testing.T) {
// 	t.Log("插入测试投票...")

// 	// 获取用户ID
// 	var users []model.User
// 	if err := mysql.Conn.Find(&users).Error; err != nil {
// 		t.Fatalf("获取用户失败: %s", err.Error())
// 	}

// 	if len(users) == 0 {
// 		t.Skip("没有用户数据，无法创建投票")
// 		return
// 	}

// 	votes := []model.Vote{
// 		{
// 			Title:       "你最喜欢的编程语言是什么？",
// 			Type:        0,     // 单选
// 			Status:      0,     // 正常
// 			Time:        86400, // 24小时
// 			UserId:      users[0].Id,
// 			CreatedTime: time.Now(),
// 			UpdatedTime: time.Now(),
// 		},
// 		{
// 			Title:       "你使用过哪些开发工具？",
// 			Type:        1,      // 多选
// 			Status:      0,      // 正常
// 			Time:        172800, // 48小时
// 			UserId:      users[0].Id,
// 			CreatedTime: time.Now(),
// 			UpdatedTime: time.Now(),
// 		},
// 		{
// 			Title:       "你最喜欢的操作系统是什么？",
// 			Type:        0,      // 单选
// 			Status:      0,      // 正常
// 			Time:        259200, // 72小时
// 			UserId:      users[1].Id,
// 			CreatedTime: time.Now(),
// 			UpdatedTime: time.Now(),
// 		},
// 		{
// 			Title:       "你希望学习哪些技术？",
// 			Type:        1,    // 多选
// 			Status:      1,    // 已结束
// 			Time:        3600, // 1小时
// 			UserId:      users[1].Id,
// 			CreatedTime: time.Now().Add(-24 * time.Hour), // 24小时前创建
// 			UpdatedTime: time.Now().Add(-23 * time.Hour), // 23小时前更新
// 		},
// 		{
// 			Title:       "你每天编程多长时间？",
// 			Type:        0,      // 单选
// 			Status:      0,      // 正常
// 			Time:        604800, // 7天
// 			UserId:      users[2].Id,
// 			CreatedTime: time.Now(),
// 			UpdatedTime: time.Now(),
// 		},
// 	}

// 	for _, vote := range votes {
// 		// 检查投票是否已存在，避免重复插入导致错误
// 		var existingVote model.Vote
// 		if err := mysql.Conn.Where("title = ?", vote.Title).First(&existingVote).Error; err == nil && existingVote.Id > 0 {
// 			t.Logf("投票 \"%s\" (ID: %d) 已存在，跳过插入", vote.Title, existingVote.Id)
// 			continue
// 		}

// 		if err := mysql.Conn.Create(&vote).Error; err != nil {
// 			t.Errorf("插入投票失败: %s", err.Error())
// 		} else {
// 			t.Logf("成功插入投票: %s (ID: %d, 类型: %s, 状态: %s)",
// 				vote.Title, vote.Id,
// 				getVoteTypeText(vote.Type),
// 				getVoteStatusText(vote.Status))
// 		}
// 	}
// }

// // testInsertTestVoteOptions 测试插入测试投票选项
// func testInsertTestVoteOptions(t *testing.T) {
// 	t.Log("插入测试投票选项...")

// 	// 获取投票ID
// 	var votes []model.Vote
// 	if err := mysql.Conn.Find(&votes).Error; err != nil {
// 		t.Fatalf("获取投票失败: %s", err.Error())
// 	}

// 	if len(votes) == 0 {
// 		t.Skip("没有投票数据，无法创建选项")
// 		return
// 	}

// 	// 为每个投票创建选项
// 	for _, vote := range votes {
// 		var options []model.VoteOpt

// 		switch vote.Id {
// 		case 1: // 编程语言投票
// 			options = []model.VoteOpt{
// 				{Name: "Go", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
// 				{Name: "Python", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
// 				{Name: "JavaScript", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
// 				{Name: "Java", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
// 				{Name: "C++", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
// 			}
// 		case 2: // 开发工具投票
// 			options = []model.VoteOpt{
// 				{Name: "VS Code", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
// 				{Name: "IntelliJ IDEA", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
// 				{Name: "Vim", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
// 				{Name: "Sublime Text", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
// 				{Name: "Atom", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
// 			}
// 		case 3: // 操作系统投票
// 			options = []model.VoteOpt{
// 				{Name: "Windows", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
// 				{Name: "macOS", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
// 				{Name: "Linux", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
// 				{Name: "Ubuntu", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
// 			}
// 		case 4: // 学习技术投票
// 			options = []model.VoteOpt{
// 				{Name: "人工智能", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
// 				{Name: "区块链", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
// 				{Name: "云计算", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
// 				{Name: "大数据", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
// 				{Name: "物联网", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
// 			}
// 		case 5: // 编程时间投票
// 			options = []model.VoteOpt{
// 				{Name: "1-2小时", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
// 				{Name: "3-4小时", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
// 				{Name: "5-6小时", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
// 				{Name: "7-8小时", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
// 				{Name: "8小时以上", VoteId: vote.Id, Count: 0, CreatedTime: time.Now(), UpdatedTime: time.Now()},
// 			}
// 		}

// 		for _, option := range options {
// 			// 检查选项是否已存在，避免重复插入导致错误
// 			var existingOpt model.VoteOpt
// 			if err := mysql.Conn.Where("name = ? AND vote_id = ?", option.Name, option.VoteId).First(&existingOpt).Error; err == nil && existingOpt.Id > 0 {
// 				t.Logf("选项 \"%s\" (投票ID: %d) 已存在，跳过插入", option.Name, option.VoteId)
// 				continue
// 			}
// 			if err := mysql.Conn.Create(&option).Error; err != nil {
// 				t.Errorf("插入投票选项失败: %s", err.Error())
// 			} else {
// 				t.Logf("成功插入选项: %s (投票ID: %d, 选项ID: %d)",
// 					option.Name, option.VoteId, option.Id)
// 			}
// 		}
// 	}
// }

// // testInsertTestVoteRecords 测试插入一些测试投票记录
// func testInsertTestVoteRecords(t *testing.T) {
// 	t.Log("插入测试投票记录...")

// 	// 获取用户和投票选项
// 	var users []model.User
// 	var voteOpts []model.VoteOpt

// 	if err := mysql.Conn.Find(&users).Error; err != nil {
// 		t.Fatalf("获取用户失败: %s", err.Error())
// 	}

// 	if err := mysql.Conn.Find(&voteOpts).Error; err != nil {
// 		t.Fatalf("获取投票选项失败: %s", err.Error())
// 	}

// 	if len(users) == 0 || len(voteOpts) == 0 {
// 		t.Skip("没有用户或投票选项数据，无法创建投票记录")
// 		return
// 	}

// 	// 按投票ID分组选项
// 	voteOptionsMap := make(map[int64][]model.VoteOpt)
// 	for _, opt := range voteOpts {
// 		voteOptionsMap[opt.VoteId] = append(voteOptionsMap[opt.VoteId], opt)
// 	}

// 	// 检查用户数量是否足够
// 	if len(users) < 3 {
// 		t.Skipf("用户数量不足（当前: %d，需要: 3），跳过投票记录创建", len(users))
// 		return
// 	}

// 	// 创建一些测试投票记录
// 	testVotes := []struct {
// 		userId     int64
// 		voteId     int64
// 		optionName string
// 	}{
// 		{users[0].Id, 1, "Go"},            // 第一个用户投票给 Go
// 		{users[1].Id, 1, "Python"},        // 第二个用户投票给 Python
// 		{users[2].Id, 1, "JavaScript"},    // 第三个用户投票给 JavaScript
// 		{users[0].Id, 2, "VS Code"},       // 第一个用户选择 VS Code
// 		{users[0].Id, 2, "IntelliJ IDEA"}, // 第一个用户选择 IntelliJ IDEA
// 		{users[1].Id, 2, "VS Code"},       // 第二个用户选择 VS Code
// 		{users[2].Id, 3, "Linux"},         // 第三个用户选择 Linux
// 		{users[0].Id, 4, "人工智能"},          // 第一个用户选择 人工智能
// 		{users[0].Id, 4, "区块链"},           // 第一个用户选择 区块链
// 		{users[1].Id, 4, "云计算"},           // 第二个用户选择 云计算
// 		{users[2].Id, 5, "3-4小时"},         // 第三个用户选择 3-4小时
// 	}

// 	for _, testVote := range testVotes {
// 		// 检查用户索引是否有效
// 		if testVote.userId <= 0 {
// 			t.Logf("无效的用户ID: %d，跳过", testVote.userId)
// 			continue
// 		}

// 		// 查找对应的投票选项
// 		var targetOpt *model.VoteOpt
// 		if options, exists := voteOptionsMap[testVote.voteId]; exists {
// 			for _, opt := range options {
// 				if opt.Name == testVote.optionName {
// 					targetOpt = &opt
// 					break
// 				}
// 			}
// 		}

// 		if targetOpt == nil {
// 			t.Logf("未找到投票 %d 的选项 \"%s\"，跳过", testVote.voteId, testVote.optionName)
// 			continue
// 		}

// 		// 检查投票记录是否已存在，避免重复插入导致错误
// 		var existingRecord model.VoteOptUser
// 		if err := mysql.Conn.Where("user_id = ? AND vote_id = ? AND vote_opt_id = ?",
// 			testVote.userId, testVote.voteId, targetOpt.Id).First(&existingRecord).Error; err == nil && existingRecord.Id > 0 {
// 			t.Logf("投票记录已存在: 用户%d 投票给选项%d (投票%d)，跳过插入",
// 				testVote.userId, targetOpt.Id, testVote.voteId)
// 			continue
// 		}

// 		// 创建投票记录
// 		voteRecord := model.VoteOptUser{
// 			UserId:      testVote.userId,
// 			VoteId:      testVote.voteId,
// 			VoteOptId:   targetOpt.Id,
// 			CreatedTime: time.Now(),
// 		}

// 		if err := mysql.Conn.Create(&voteRecord).Error; err != nil {
// 			t.Errorf("插入投票记录失败: %s", err.Error())
// 		} else {
// 			t.Logf("成功插入投票记录: 用户%d 投票给选项\"%s\" (投票%d)",
// 				testVote.userId, testVote.optionName, testVote.voteId)
// 		}

// 		// 更新选项计数 (注意：这里的计数更新应该在提交投票时由业务逻辑完成，这里只是为了测试数据完整性)
// 		if err := mysql.Conn.Model(&model.VoteOpt{}).Where("id = ?", targetOpt.Id).
// 			Update("count", mysql.Conn.Raw("count + 1")).Error; err != nil {
// 			t.Errorf("更新选项计数失败: %s", err.Error())
// 		}
// 	}
// }

// // TestPrintTestData 测试打印测试数据统计
// func TestPrintTestData(t *testing.T) {
// 	// 初始化数据库连接
// 	mysql.NewMysql()
// 	defer mysql.Close()

// 	t.Log("\n=== 测试数据统计 ===")

// 	// 统计用户
// 	var userCount int64
// 	mysql.Conn.Model(&model.User{}).Count(&userCount)
// 	t.Logf("用户总数: %d", userCount)

// 	// 统计投票
// 	var voteCount int64
// 	mysql.Conn.Model(&model.Vote{}).Count(&voteCount)
// 	t.Logf("投票总数: %d", voteCount)

// 	// 统计投票选项
// 	var voteOptCount int64
// 	mysql.Conn.Model(&model.VoteOpt{}).Count(&voteOptCount)
// 	t.Logf("投票选项总数: %d", voteOptCount)

// 	// 统计投票记录
// 	var voteRecordCount int64
// 	mysql.Conn.Model(&model.VoteOptUser{}).Count(&voteRecordCount)
// 	t.Logf("投票记录总数: %d", voteRecordCount)

// 	t.Log("=== 数据统计完成 ===")
// }

// // TestDataIntegrity 测试数据完整性
// func TestDataIntegrity(t *testing.T) {
// 	// 初始化数据库连接
// 	mysql.NewMysql()
// 	defer mysql.Close()

// 	t.Log("开始数据完整性测试...")

// 	// 检查用户数据完整性
// 	var users []model.User
// 	if err := mysql.Conn.Find(&users).Error; err != nil {
// 		t.Fatalf("获取用户数据失败: %s", err)
// 	}

// 	for _, user := range users {
// 		if user.Name == "" {
// 			t.Errorf("用户ID %d 的用户名为空", user.Id)
// 		}
// 		if user.Password == "" {
// 			t.Errorf("用户 %s 的密码为空", user.Name)
// 		}
// 		if user.Uuid == 0 {
// 			t.Errorf("用户 %s 的UUID为0", user.Name)
// 		}
// 	}

// 	// 检查投票数据完整性
// 	var votes []model.Vote
// 	if err := mysql.Conn.Find(&votes).Error; err != nil {
// 		t.Fatalf("获取投票数据失败: %s", err)
// 	}

// 	for _, vote := range votes {
// 		if vote.Title == "" {
// 			t.Errorf("投票ID %d 的标题为空", vote.Id)
// 		}
// 		if vote.UserId == 0 {
// 			t.Errorf("投票 %s 的创建者ID为0", vote.Title)
// 		}
// 	}

// 	// 检查投票选项数据完整性
// 	var voteOpts []model.VoteOpt
// 	if err := mysql.Conn.Find(&voteOpts).Error; err != nil {
// 		t.Fatalf("获取投票选项数据失败: %s", err)
// 	}

// 	for _, opt := range voteOpts {
// 		if opt.Name == "" {
// 			t.Errorf("投票选项ID %d 的名称为空", opt.Id)
// 		}
// 		if opt.VoteId == 0 {
// 			t.Errorf("投票选项 %s 的投票ID为0", opt.Name)
// 		}
// 	}

// 	t.Log("数据完整性测试完成")
// }

// // getVoteTypeText 获取投票类型文本
// func getVoteTypeText(voteType int32) string {
// 	if voteType == 0 {
// 		return "单选"
// 	}
// 	return "多选"
// }

// // getVoteStatusText 获取投票状态文本
// func getVoteStatusText(voteStatus int32) string {
// 	if voteStatus == 0 {
// 		return "正常"
// 	}
// 	return "已结束"
// }

// // 保留原有的函数用于向后兼容
// func InsertTestData() {
// 	fmt.Println("开始插入测试数据...")
// 	// 这里可以调用测试函数，但为了保持兼容性，保留原有逻辑
// }

// func PrintTestData() {
// 	fmt.Println("\n=== 测试数据统计 ===")
// 	// 这里可以调用测试函数，但为了保持兼容性，保留原有逻辑
// }
