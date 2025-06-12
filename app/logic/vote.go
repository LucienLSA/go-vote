package logic

import (
	"govote/app/model"
	"govote/app/tools/e"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

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

	opt := make([]model.VoteOpt, 0)
	for _, v := range optStr {
		opt = append(opt, model.VoteOpt{
			Name:        v,
			CreatedTime: time.Now(),
		})
	}

	if err := model.AddVote(vote, opt); err != nil {
		context.JSON(http.StatusOK, e.ECode{
			Code:    10006,
			Message: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, e.OK)
}

func UpdateVote(context *gin.Context) {

}

// DelVote 删除一个投票
func DelVote(context *gin.Context) {
	var id int64
	idStr := context.Query("id")
	id, _ = strconv.ParseInt(idStr, 10, 64)
	if ok := model.DelVote(id); !ok {
		context.JSON(http.StatusOK, e.ECode{
			Code: 10006,
		})
		return
	}

	context.JSON(http.StatusOK, e.OK)
}
