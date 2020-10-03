package utils

import (
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// GetModuleLogger - 提供一个为 Module 使用的 logrus.Entry
// 包含 logrus.Fields
func GetModuleLogger(name string) *logrus.Entry {
	return logrus.WithField("module", name)
}

// WriteLogToFS 将日志转储至文件
// 请务必在 init() 阶段调用此函数
// 否则会出现日志缺失
// 日志存储位置 ./logs
func WriteLogToFS() {
	WriteLogToPath("logs")
}

// WriteLogToPath 将日志转储至文件
// 请务必在 init() 阶段调用此函数
// 否则会出现日志缺失
// 日志存储位置 p
func WriteLogToPath(p string) {
	writer, err := rotatelogs.New(
		path.Join(p, "%Y-%m-%d.log"),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		logrus.WithError(err).Error("unable to write logs")
		return
	}
	logrus.AddHook(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.InfoLevel:  writer,
			logrus.ErrorLevel: writer,
			logrus.FatalLevel: writer,
		}, &logrus.JSONFormatter{},
	))
}
