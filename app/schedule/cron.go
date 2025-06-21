package schedule

import (
	"context"
	"govote/app/db/mysql"
	"log"
	"time"
)

func Start(ctx context.Context, interval time.Duration) {
	go EndVoteCorn(ctx, interval)
}

func EndVoteCorn(ctx context.Context, interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("EndVote定时任务已停止")
			return
		case <-t.C:
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("EndVote定时任务发生panic: %v", r)
					}
				}()
				log.Println("EndVote启动")
				mysql.EndVote(ctx)
				log.Println("EndVote运行完毕")
			}()
		}
	}
}
