package logic

import (
	"govote/app/model"
	"govote/app/tools/e"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// AddVote godoc
// @Summary      新增投票
// @Description  新增投票
// @Tags         vote
// @Accept       json
// @Produce      json
// @Param        title   query     string  true	"vote string"
// @Success      200  {object}  e.ECode
// @Router       /vote/add [post]
func AddVote(context *gin.Context) {
	idStr := context.Query("title")
	optStr, _ := context.GetPostFormArray("opt_name[]")
	//构建结构体
	vote := model.Vote{
		Title:       idStr,
		Type:        0,
		Status:      0,
		CreatedTime: time.Now(),
	}
	if vote.Title == "" {
		context.JSON(http.StatusBadRequest, e.ParamErr)
		return
	}
	// 幂等性，在添加投票记录前查询是否存在
	oldVote := model.GetVoteByName(idStr)
	if oldVote.Id > 0 {
		context.JSON(http.StatusOK, e.VoteRepeatErr)
		return
	}
	opt := make([]model.VoteOpt, 0)
	for _, v := range optStr {
		opt = append(opt, model.VoteOpt{
			Name:        v,
			CreatedTime: time.Now(),
		})
	}

	if err := model.AddVote(vote, opt); err != nil {
		context.JSON(http.StatusOK, e.ServerErr)
		return
	}

	context.JSON(http.StatusCreated, e.OK)
}

// UpdateVote godoc
// @Summary      更新投票
// @Description  更新投票
// @Tags         vote
// @Accept       json
// @Produce      json
// @Param        title   query     string  true	"vote string"
// @Success      200  {object}  e.ECode
// @Router       /vote/update [post]
func UpdateVote(context *gin.Context) {
	idStr := context.Query("title")
	optStr, _ := context.GetPostFormArray("opt_name[]")
	//构建结构体
	vote := model.Vote{
		Title:       idStr,
		Type:        0,
		Status:      0,
		CreatedTime: time.Now(),
	}

	opt := make([]model.VoteOpt, 0)
	for _, v := range optStr {
		opt = append(opt, model.VoteOpt{
			Name:        v,
			CreatedTime: time.Now(),
		})
	}
	if err := model.UpdateVote(vote, opt); err != nil {
		context.JSON(http.StatusOK, e.ServerErr)
		return
	}

	context.JSON(http.StatusCreated, e.OK)
}

// DelVote godoc
// @Summary      删除投票
// @Description  删除投票
// @Tags         vote
// @Accept       json
// @Produce      json
// @Param        title   query     string  true	"vote string"
// @Success      200  {object}  e.ECode
// @Router       /vote/del [post]
func DelVote(context *gin.Context) {
	var id int64
	idStr := context.Query("id")
	id, _ = strconv.ParseInt(idStr, 10, 64)
	// vote := model.GetVote(id)
	voteInfo := model.GetVoteCache(context, id)
	if voteInfo.Vote.Id < 1 {
		context.JSON(http.StatusNoContent, e.OK)
		return
	}
	if ok := model.DelVote(id); !ok {
		context.JSON(http.StatusOK, e.ServerErr)
		return
	}
	context.JSON(http.StatusNoContent, e.OK)
}
