package model

import "time"

type User struct {
	Id          int64  `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	Name        string `gorm:"column:name;default:NULL"`
	Password    string `gorm:"column:password;default:NULL"`
	CreatedTime string `gorm:"column:created_time;default:NULL"`
	UpdatedTime string `gorm:"column:updated_time;default:NULL"`
}

// TableName 表名
func (u *User) TableName() string {
	return "user"
}

type Vote struct {
	Id          int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"id"`
	Title       string    `gorm:"column:title;type:varchar(255)" json:"title"`
	Type        int       `gorm:"column:type;type:int(11);comment:0单选1多选" json:"type"`
	Status      int       `gorm:"column:status;type:int(11);comment:0正常1超时" json:"status"`
	Time        int64     `gorm:"column:time;type:bigint(20);comment:有效时长" json:"time"`
	UserId      int64     `gorm:"column:user_id;type:bigint(20);comment:创建人" json:"user_id"`
	CreatedTime time.Time `gorm:"column:created_time;type:datetime" json:"created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time;type:datetime" json:"updated_time"`
}

func (m *Vote) TableName() string {
	return "vote"
}

type VoteOptUser struct {
	Id          int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"id"`
	UserId      int64     `gorm:"column:user_id;type:bigint(20)" json:"user_id"`
	VoteId      int64     `gorm:"column:vote_id;type:bigint(20)" json:"vote_id"`
	VoteOptId   int64     `gorm:"column:vote_opt_id;type:bigint(20)" json:"vote_opt_id"`
	CreatedTime time.Time `gorm:"column:created_time;type:datetime" json:"created_time"`
}

func (m *VoteOptUser) TableName() string {
	return "vote_opt_user"
}

type VoteOpt struct {
	Id          int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"id"`
	Name        string    `gorm:"column:name;type:varchar(255)" json:"name"`
	VoteId      int64     `gorm:"column:vote_id;type:bigint(20)" json:"vote_id"`
	Count       int       `gorm:"column:count;type:int(11)" json:"count"`
	CreatedTime time.Time `gorm:"column:created_time;type:datetime" json:"created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time;type:datetime" json:"updated_time"`
}

func (m *VoteOpt) TableName() string {
	return "vote_opt"
}
