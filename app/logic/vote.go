package logic

import (
	"govote/app/model"
	"govote/app/types"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetVoteDetail 获取投票详情页面
func GetVoteDetail(context *gin.Context) {
	// 获取投票ID
	voteIDStr := context.Param("id")
	voteID, err := strconv.ParseInt(voteIDStr, 10, 64)
	if err != nil {
		context.HTML(http.StatusBadRequest, "vote.html", gin.H{
			"error": "无效的投票ID",
		})
		return
	}

	// 获取投票信息
	voteInfo := &types.Vote{Id: voteID}
	vote := model.GetVote(voteInfo)

	// 检查投票是否存在
	if vote.Id == 0 {
		context.HTML(http.StatusNotFound, "vote.html", gin.H{
			"error": "投票不存在",
		})
		return
	}

	// 获取投票选项
	voteOpts := model.GetVoteOpts(voteID)
	if len(voteIDStr) == 0 {
		context.HTML(http.StatusInternalServerError, "vote.html", gin.H{
			"error": "获取投票选项失败",
			"vote":  vote,
		})
		return
	}

	// 创建模板数据
	data := gin.H{
		"vote":     vote,
		"voteOpts": voteOpts,
	}

	context.HTML(http.StatusOK, "vote.html", data)
}

// SubmitVote 提交投票
func SubmitVote(context *gin.Context) {
	// 获取当前登录用户
	name, err := context.Cookie("name")
	if err != nil || name == "" {
		context.HTML(http.StatusUnauthorized, "vote.html", gin.H{
			"error": "请先登录",
		})
		return
	}

	// 获取投票ID
	voteIDStr := context.PostForm("vote_id")
	voteID, err := strconv.ParseInt(voteIDStr, 10, 64)
	if err != nil {
		context.HTML(http.StatusBadRequest, "vote.html", gin.H{
			"error": "无效的投票ID",
		})
		return
	}

	// 获取投票信息
	voteInfo := &types.Vote{Id: voteID}
	vote := model.GetVote(voteInfo)

	// 检查投票是否存在
	if vote.Id == 0 {
		context.HTML(http.StatusNotFound, "vote.html", gin.H{
			"error": "投票不存在",
		})
		return
	}

	// 检查投票状态
	if vote.Status != 0 {
		context.HTML(http.StatusBadRequest, "vote.html", gin.H{
			"error": "投票已结束，无法参与",
			"vote":  vote,
		})
		return
	}

	// 获取用户选择的选项
	var selectedOpts []string
	if vote.Type == 0 {
		// 单选
		opt := context.PostForm("opt")
		if opt == "" {
			context.HTML(http.StatusBadRequest, "vote.html", gin.H{
				"error": "请选择一个选项",
				"vote":  vote,
			})
			return
		}
		selectedOpts = append(selectedOpts, opt)
	} else {
		// 多选
		opts := context.PostFormArray("opt[]")
		if len(opts) == 0 {
			context.HTML(http.StatusBadRequest, "vote.html", gin.H{
				"error": "请至少选择一个选项",
				"vote":  vote,
			})
			return
		}
		selectedOpts = opts
	}

	// 获取用户ID
	user := model.GetUserByName(name)
	if user.Id == 0 {
		context.HTML(http.StatusInternalServerError, "vote.html", gin.H{
			"error": "用户信息获取失败",
			"vote":  vote,
		})
		return
	}

	// 检查用户是否已经投票
	for _, optIDStr := range selectedOpts {
		optID, err := strconv.ParseInt(optIDStr, 10, 64)
		if err != nil {
			continue
		}

		existingRecord := model.CheckUserVote(user.Id, voteID, optID)
		if existingRecord.Id != 0 {
			context.HTML(http.StatusBadRequest, "vote.html", gin.H{
				"error": "您已经参与过此投票",
				"vote":  vote,
			})
			return
		}
	}

	// 提交投票
	for _, optIDStr := range selectedOpts {
		optID, err := strconv.ParseInt(optIDStr, 10, 64)
		if err != nil {
			continue
		}

		// 创建投票记录
		err = model.CreateVote(user.Id, voteID, optID)
		if err != nil {
			context.HTML(http.StatusInternalServerError, "vote.html", gin.H{
				"error": "投票提交失败",
				"vote":  vote,
			})
			return
		}

		// 更新选项投票数
		err = model.DB.Model(&model.VoteOpt{}).Where("id = ?", optID).
			Update("count", model.DB.Raw("count + 1")).Error
		if err != nil {
			context.HTML(http.StatusInternalServerError, "vote.html", gin.H{
				"error": "投票统计更新失败",
				"vote":  vote,
			})
			return
		}
	}

	// 获取更新后的投票选项
	var voteOpts []model.VoteOpt
	err = model.DB.Where("vote_id = ?", voteID).Find(&voteOpts).Error
	if err != nil {
		context.HTML(http.StatusInternalServerError, "vote.html", gin.H{
			"error": "获取投票选项失败",
			"vote":  vote,
		})
		return
	}

	// 显示成功消息
	context.HTML(http.StatusOK, "vote.html", gin.H{
		"success":  "投票提交成功！",
		"vote":     vote,
		"voteOpts": voteOpts,
	})
}

// GetVoteResult 获取投票结果
func GetVoteResult(context *gin.Context) {
	// 获取投票ID
	voteIDStr := context.Param("id")
	voteID, err := strconv.ParseInt(voteIDStr, 10, 64)
	if err != nil {
		context.HTML(http.StatusBadRequest, "vote.html", gin.H{
			"error": "无效的投票ID",
		})
		return
	}

	// 获取投票信息
	voteInfo := &types.Vote{Id: voteID}
	vote := model.GetVote(voteInfo)

	// 检查投票是否存在
	if vote.Id == 0 {
		context.HTML(http.StatusNotFound, "vote.html", gin.H{
			"error": "投票不存在",
		})
		return
	}

	// 获取投票选项（按投票数排序）
	var voteOpts []model.VoteOpt
	err = model.DB.Where("vote_id = ?", voteID).Order("count desc").Find(&voteOpts).Error
	if err != nil {
		context.HTML(http.StatusInternalServerError, "vote.html", gin.H{
			"error": "获取投票结果失败",
			"vote":  vote,
		})
		return
	}

	// 计算总投票数
	totalVotes := 0
	for _, opt := range voteOpts {
		totalVotes += opt.Count
	}

	// 创建模板数据
	data := gin.H{
		"vote":       vote,
		"voteOpts":   voteOpts,
		"totalVotes": totalVotes,
		"isResult":   true,
	}

	context.HTML(http.StatusOK, "vote.html", data)
}
