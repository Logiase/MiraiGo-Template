package utils

import (
	"fmt"
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type LogArg uint32

const (
	LogFatalLevel        = LogArg(logrus.FatalLevel)
	LogErrorLevel        = LogArg(logrus.ErrorLevel)
	LogWarnLevel         = LogArg(logrus.WarnLevel)
	LogInfoLevel         = LogArg(logrus.InfoLevel)
	LogDebugLevel        = LogArg(logrus.DebugLevel)
	LogTraceLevel        = LogArg(logrus.TraceLevel)
	logLevelMask  LogArg = 0x07
	LogWithStack  LogArg = 0x08
)

var logWithStack = false

// GetModuleLogger - 提供一个为 Module 使用的 logrus.Entry
// 包含 logrus.Fields
func GetModuleLogger(name string) logrus.FieldLogger {
	if logWithStack {
		return &errorEntryWithStack{logrus.WithField("module", name)}
	} else {
		return logrus.WithField("module", name)
	}
}

// WriteLogToFS 将日志转储至文件
// 请务必在 init() 阶段调用此函数
// 否则会出现日志缺失
// 日志存储位置 ./logs
func WriteLogToFS(args ...LogArg) {
	WriteLogToPath("logs", args...)
}

// WriteLogToPath 将日志转储至文件
// 请务必在 init() 阶段调用此函数
// 否则会出现日志缺失
// 日志存储位置 p
func WriteLogToPath(p string, args ...LogArg) {
	writer, err := rotatelogs.New(
		path.Join(p, "%Y-%m-%d.log"),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		logrus.WithError(err).Error("unable to write logs")
		return
	}

	var arg LogArg
	for _, a := range args {
		arg |= a
	}
	if arg&LogWithStack == LogWithStack {
		logWithStack = true
	}
	logLevel := arg & logLevelMask
	if logLevel == 0 {
		logLevel = LogInfoLevel
	}

	writerMap := make(lfshook.WriterMap)
	switch {
	case logLevel > LogFatalLevel:
		writerMap[logrus.FatalLevel] = writer
		fallthrough
	case logLevel > LogErrorLevel:
		writerMap[logrus.ErrorLevel] = writer
		fallthrough
	case logLevel > LogWarnLevel:
		writerMap[logrus.WarnLevel] = writer
		fallthrough
	case logLevel > LogInfoLevel:
		writerMap[logrus.InfoLevel] = writer
		fallthrough
	case logLevel > LogDebugLevel:
		writerMap[logrus.DebugLevel] = writer
		fallthrough
	case logLevel > LogTraceLevel:
		writerMap[logrus.TraceLevel] = writer
	}
	if logWithStack {
		logrus.AddHook(lfshook.NewHook(writerMap, &logrus.TextFormatter{DisableQuote: true}))
	} else {
		logrus.AddHook(lfshook.NewHook(writerMap, &logrus.JSONFormatter{}))
	}
}

type errorEntryWithStack struct {
	*logrus.Entry
}

func (e *errorEntryWithStack) WithError(err error) *logrus.Entry {
	return e.Entry.WithError(fmt.Errorf("%+v", err))
}
