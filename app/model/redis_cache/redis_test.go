package redis_cache

import (
	"context"
	"fmt"
	"govote/app/model/mysql"
	"govote/app/tools/log"
	"testing"
)

func TestGetHistoryVote(t *testing.T) {
	log.NewLogger()
	mysql.NewMysql()
	NewRedis() // 初始化Redis连接
	//测试用例
	r, _ := GetVoteUserHistory(context.Background(), 1, 1)
	fmt.Printf("ret:%+v", r)
	mysql.Close()
}
