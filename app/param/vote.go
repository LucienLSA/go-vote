package param

// 投票记录 参数
type VoteData struct {
	Id int64 `json:"id" form:"id"`
}

// 投票信息 参数
type VoteInfoData struct {
	UserID int64   `json:"user_id" form:"user_id"`
	VoteID int64   `json:"vote_id" form:"vote_id"`
	OptIDs []int64 `json:"opt_id" form:"opt_id"`
}
