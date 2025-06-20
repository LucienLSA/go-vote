package schedule

import (
	"fmt"
	"govote/app/model/mysql"
	"time"
)

func Start() {
	go EndVote()
}

func EndVote() {
	t := time.NewTicker(5 * time.Second)
	defer func() {
		t.Stop()
	}()

	for {
		select {
		case <-t.C:
			fmt.Println("EndVote启动")
			// 执行的函数
			mysql.EndVote()
			fmt.Println("EndVote运行完毕")
		}
	}
}
