package logic

import (
	"govote/app/model"
	"govote/app/param"
	"govote/app/tools/e"
	"govote/app/tools/session"
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
// @Param 		 id    body    param.VoteData   true  "vote param.VoteData"
// @Success      200  {object}  e.ECode
// @Router       /vote [get]
func GetVoteInfo(context *gin.Context) {
	var voteData param.VoteData
	idStr := context.Query("id")
	voteData.Id, _ = strconv.ParseInt(idStr, 10, 64)
	// ret := model.GetVote(voteData.Id)
	ret := model.GetVoteCache(context, voteData.Id)
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
// @Param 		 Id    body     param.VoteInfoData   true  "user  param.VoteInfoData"
// @Success      200  {object}  e.ECode
// @Router       /vote [post]
func DoVote(context *gin.Context) {
	var voteInfo param.VoteInfoData
	// 使用session机制获取用户ID
	values := session.GetSession(context)
	var userID int64
	if v, ok := values["id"]; ok {
		userID = v.(int64)
	}
	voteIdStr, _ := context.GetPostForm("vote_id")
	optStr, _ := context.GetPostFormArray("opt[]")

	voteInfo.UserID = userID
	voteInfo.VoteID, _ = strconv.ParseInt(voteIdStr, 10, 64)

	//查询是否投过票了
	voteUser, err := model.GetVoteUserHistory(context, voteInfo.UserID, voteInfo.VoteID)
	if len(voteUser) > 0 || err != nil {
		context.JSON(http.StatusOK, e.VoteRepeatErr)
		return
	}

	voteInfo.Opt = make([]int64, 0)
	for _, v := range optStr {
		optId, _ := strconv.ParseInt(v, 10, 64)
		voteInfo.Opt = append(voteInfo.Opt, optId)
	}
	// 执行投票
	ok := model.DoVoteV2(voteInfo.UserID, voteInfo.VoteID, voteInfo.Opt)
	if !ok {
		context.JSON(http.StatusOK, e.ServerErr)
		return
	}
	// 投票完成删除缓存，设置为过期
	err = model.CleanVote(context, voteInfo.UserID)
	if err != nil {
		context.JSON(http.StatusOK, e.ServerErr)
		return
	}

	context.JSON(http.StatusOK, e.ECode{
		Code:    0,
		Message: "投票成功",
	})
}
