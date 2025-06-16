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

// GetVoteInfo godoc
// @Summary      获取投票信息
// @Description  获取投票信息
// @Tags         vote
// @Accept       json
// @Produce      json
// @Param 		 id    query    int   true  "vote Id"
// @Success      200  {object}  e.ECode
// @Router       /vote [get]
func GetVoteInfo(context *gin.Context) {
	var id int64
	idStr := context.Query("id")
	id, _ = strconv.ParseInt(idStr, 10, 64)
	ret := model.GetVote(id)
	// context.HTML(http.StatusOK, "vote.html", gin.H{"vote": ret})
	if ret.Vote.Id < 1 {
		context.JSON(http.StatusNotFound, e.NotFoundErr)
	}
	context.JSON(http.StatusOK, e.ECode{
		Data: ret,
	})
}

// DoVote godoc
// @Summary      投票
// @Description  投票
// @Tags         vote
// @Accept       json
// @Produce      json
// @Param 		 Id    query    int   true  "user Id"
// @Param 		 vote_id    query    int   true  "vote Id"
// @Param 		 opt[]    query    int   true  "vote_opt"
// @Success      200  {object}  e.ECode
// @Router       /vote [post]
func DoVote(context *gin.Context) {
	userIDStr, _ := context.Cookie("Id")
	voteIdStr, _ := context.GetPostForm("vote_id")
	optStr, _ := context.GetPostFormArray("opt[]")

	userID, _ := strconv.ParseInt(userIDStr, 10, 64)
	voteId, _ := strconv.ParseInt(voteIdStr, 10, 64)

	//查询是否投过票了
	voteUser := model.GetVoteHistory(userID, voteId)
	if len(voteUser) > 0 {
		context.JSON(http.StatusOK, e.VoteRepeatErr)
		return
	}

	opt := make([]int64, 0)
	for _, v := range optStr {
		optId, _ := strconv.ParseInt(v, 10, 64)
		opt = append(opt, optId)
	}

	model.DoVoteV2(userID, voteId, opt)
	context.JSON(http.StatusOK, e.ECode{
		Code:    0,
		Message: "投票成功",
	})
}
