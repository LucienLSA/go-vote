package types

// Vote 投票结构体
type Vote struct {
	Id    int64  `json:"id" form:"id"`
	Title string `json:"title" form:"title"`
}

type VoteOpts struct {
	Id     int64  ` json:"id" form:"id"`
	Name   string ` json:"name" form:"name"`
	VoteId int64  ` json:"vote_id" form:"vote_id"`
}

type VoteOptUser struct {
	Id        int64 ` json:"id" form:"id"`
	UserId    int64 `json:"user_id" form:"user_id"`
	VoteId    int64 ` json:"vote_id" form:"vote_id"`
	VoteOptId int64 ` json:"vote_opt_id" form:"vote_opt_id"`
}
