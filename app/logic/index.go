package logic

import (
	"govote/app/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(context *gin.Context) {
	// 获取当前登录用户
	name, err := context.Cookie("name")
	userName := ""
	if err == nil && name != "" {
		userName = name
	}
	// 从数据库获取真实的投票数据
	dbVotes := model.GetVotes()

	// 转换为模板需要的格式
	votes := make([]map[string]interface{}, 0, len(dbVotes))
	normalCount := 0
	expiredCount := 0
	singleCount := 0
	multiCount := 0

	for _, vote := range dbVotes {
		voteData := map[string]interface{}{
			"Id":     vote.Id,
			"Title":  vote.Title,
			"Type":   vote.Type,
			"Status": vote.Status,
			"Time":   vote.Time,
			"UserId": vote.UserId,
		}
		votes = append(votes, voteData)

		// 统计数据
		if vote.Status == 0 {
			normalCount++
		} else {
			expiredCount++
		}

		if vote.Type == 0 {
			singleCount++
		} else {
			multiCount++
		}
	}

	// 创建模板数据
	data := gin.H{
		"vote":         votes,
		"userName":     userName,
		"isLogin":      userName != "",
		"normalCount":  normalCount,
		"expiredCount": expiredCount,
		"singleCount":  singleCount,
		"multiCount":   multiCount,
	}

	context.HTML(http.StatusOK, "index.html", data)
}
