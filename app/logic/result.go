package logic

import (
	"govote/app/db/redis_cache"
	"govote/app/tools/e"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ResultData 新定义返回结构
type ResultData struct {
	Title string
	Count int64
	Opt   []*ResultVoteOpt
}

type ResultVoteOpt struct {
	Name  string
	Count int64
}

func ResultInfo(context *gin.Context) {
	context.HTML(http.StatusOK, "result.html", nil)
}

// ResultVote 返回一个投票结果
func ResultVote(context *gin.Context) {
	var id int64
	idStr := context.Query("id")
	id, _ = strconv.ParseInt(idStr, 10, 64)
	// ret := model.GetVote(id)
	ret := redis_cache.GetVoteCache(context, id)
	data := ResultData{
		Title: ret.Vote.Title,
	}

	for _, v := range ret.Opt {
		data.Count = data.Count + v.Count
		tmp := ResultVoteOpt{
			Name:  v.Name,
			Count: v.Count,
		}
		data.Opt = append(data.Opt, &tmp)
	}

	context.JSON(http.StatusOK, e.ECode{
		Data: data,
	})
}
