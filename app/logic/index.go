package logic

import (
	"govote/app/model"
	"govote/app/tools/e"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Index(context *gin.Context) {
	ret := model.GetVotes()
	context.HTML(http.StatusOK, "index.html", gin.H{"vote": ret})
}

func GetVotes(context *gin.Context) {
	ret := model.GetVotes()
	context.JSON(http.StatusOK, e.ECode{
		Data: ret,
	})
}

func GetVoteInfo(context *gin.Context) {
	var id int64
	idStr := context.Query("id")
	id, _ = strconv.ParseInt(idStr, 10, 64)
	ret := model.GetVote(id)
	// context.HTML(http.StatusOK, "vote.html", gin.H{"vote": ret})
	context.JSON(http.StatusOK, e.ECode{
		Data: ret,
	})
}

func DoVote(context *gin.Context) {
	userIDStr, _ := context.Cookie("Id")
	voteIdStr, _ := context.GetPostForm("vote_id")
	optStr, _ := context.GetPostFormArray("opt[]")

	userID, _ := strconv.ParseInt(userIDStr, 10, 64)
	voteId, _ := strconv.ParseInt(voteIdStr, 10, 64)

	//查询是否投过票了
	voteUser := model.GetVoteHistory(userID, voteId)
	if len(voteUser) > 0 {
		context.JSON(http.StatusOK, e.ECode{
			Code:    10010,
			Message: "您已投过票了",
		})
		return
	}

	opt := make([]int64, 0)
	for _, v := range optStr {
		optId, _ := strconv.ParseInt(v, 10, 64)
		opt = append(opt, optId)
	}

	model.DoVoteV1(userID, voteId, opt)
	context.JSON(http.StatusOK, e.ECode{
		Message: "投票完成",
	})
}
