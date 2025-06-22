package log

import (
	"io"
	"os"
	"path/filepath"

	"govote/app/config"

	"github.com/sirupsen/logrus"
)

// var L *logrus.Entry
var L *logrus.Logger

// func NewLogger() {
// 	LStore := logrus.New()
// 	// 设置日志级别
// 	LStore.SetLevel(logrus.DebugLevel)
// 	// 将日志输出到控制台和文件中
// 	w1 := os.Stdout
// 	w2, _ := os.OpenFile("./vote.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
// 	LStore.SetOutput(io.MultiWriter(w1, w2))

// 	// JSON格式日志输出
// 	L = LStore.WithFields(logrus.Fields{
// 		"name": "lucien",
// 		"app":  "voteV2",
// 	})
// }

func NewLogger() {
	L = logrus.New()

	// 从配置文件获取日志级别
	level, err := logrus.ParseLevel(config.Conf.LogConfig.Level)
	if err != nil {
		// 如果解析失败，默认使用Info级别
		level = logrus.InfoLevel
	}
	L.SetLevel(level)

	// 确保日志目录存在
	logDir := config.Conf.LogConfig.FilePath
	if err := os.MkdirAll(logDir, 0755); err != nil {
		L.Errorf("创建日志目录失败: %v", err)
	}

	// 构建完整的日志文件路径
	logFilePath := filepath.Join(logDir, config.Conf.LogConfig.Filename)

	// 将日志输出到控制台和文件中
	w1 := os.Stdout
	w2, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		L.Errorf("打开日志文件失败: %v", err)
		// 如果文件打开失败，只输出到控制台
		L.SetOutput(w1)
		return
	}

	L.SetOutput(io.MultiWriter(w1, w2))
	// L.AddHook()
}
