package schedule

import (
	"fmt"
	"govote/app/model"
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
			model.EndVote()
			fmt.Println("EndVote运行完毕")
		}
	}
}
