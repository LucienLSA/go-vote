package logic

import (
	"govote/app/db/mysql"
	"govote/app/db/redis_cache"
	"govote/app/param"
	"govote/app/tools/e"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Index(context *gin.Context) {
	ret := mysql.GetVotes(context)
	context.HTML(http.StatusOK, "index.html", gin.H{"vote": ret})
}

func GetVotes(context *gin.Context) {
	ret := mysql.GetVotes(context)
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
	ret := redis_cache.GetVoteCache(context, voteData.Id)
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
	// 绑定 JSON 数据
	if err := context.ShouldBindJSON(&voteInfo); err != nil {
		context.JSON(http.StatusBadRequest, e.ParamErr)
		return
	}

	// JWT上下文获取用户id
	uid, err := GetLoginUserID(context)
	if err != nil {
		context.JSON(http.StatusOK, e.NotLogin) // NotLogin包含了UserNotLogin的错误信息
		return
	}
	voteInfo.UserID = uid

	//查询是否投过票了
	voteUser, err := redis_cache.GetVoteUserHistory(context, voteInfo.UserID, voteInfo.VoteID)
	if err != nil && err.Error() != "redis: nil" { // 忽略 "redis: nil" 错误，它表示缓存未命中
		context.JSON(http.StatusOK, e.ServerErr)
		return
	}
	if len(voteUser) > 0 {
		context.JSON(http.StatusOK, e.VoteRepeatErr)
		return
	}

	// 执行投票
	ok := mysql.DoVoteV2(context, voteInfo.UserID, voteInfo.VoteID, voteInfo.OptIDs)
	if !ok {
		context.JSON(http.StatusOK, e.ServerErr)
		return
	}
	// 投票完成删除缓存，设置为过期
	err = redis_cache.CleanVote(context, voteInfo.VoteID) // 清理投票详情缓存
	if err != nil {
		// 此处错误可以只记录日志，不一定需要返回给用户
		// log.L.Errorf("清理投票缓存失败: %v", err)
	}

	context.JSON(http.StatusOK, e.ECode{
		Code:    e.OK.Code,
		Message: "投票成功",
	})
}
