package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var L *logrus.Entry

func NewLogger() {
	LStore := logrus.New()
	// 设置日志级别
	LStore.SetLevel(logrus.DebugLevel)
	// 将日志输出到控制台和文件中
	w1 := os.Stdout
	w2, _ := os.OpenFile("./vote.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	LStore.SetOutput(io.MultiWriter(w1, w2))

	// JSON格式日志输出
	L = LStore.WithFields(logrus.Fields{
		"name": "lucien",
		"app":  "voteV2",
	})
}
