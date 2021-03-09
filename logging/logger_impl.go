package logging

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
)

// internal logger output level format
var (
	DEBUG string = color.New(color.BgBlue).SprintfFunc()("DEBUG")
	INFO  string = color.New(color.BgGreen).SprintfFunc()("INFO ")
	WARN  string = color.New(color.BgYellow).SprintfFunc()("WARN ")
	ERROR string = color.New(color.BgRed).SprintfFunc()("ERROR")
	PANIC string = color.New(color.BgCyan).SprintfFunc()("PANIC")
	FATAL string = color.New(color.BgHiRed).SprintfFunc()("FATAL")
)

const (
	// LogFormat is internal logger output format
	LogFormat = "%s ID(%s): %s"
)

// Logger is an implementation of ILogger
// using standard log
//
// Logger 是 ILogger 的内置实现
// 使用 standard log 库
type Logger struct {
	ID string

	debug bool
}

// NewLogger 创建新 Logger 实例
func NewLogger(id string) Logger {
	return Logger{
		ID:    id,
		debug: false,
	}
}

// Printf print a simple message
func (l *Logger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// Debugf print debug message
func (l *Logger) Debugf(format string, v ...interface{}) {
	if !l.debug {
		return
	}
	log.Printf(LogFormat, DEBUG, l.ID, fmt.Sprintf(format, v...))
}

// Infof print info message
func (l *Logger) Infof(format string, v ...interface{}) {
	log.Printf(LogFormat, INFO, l.ID, fmt.Sprintf(format, v...))
}

// Warnf print warn message
func (l *Logger) Warnf(format string, v ...interface{}) {
	log.Printf(LogFormat, WARN, l.ID, fmt.Sprintf(format, v...))
}

// Errorf print error message
func (l *Logger) Errorf(format string, v ...interface{}) {
	log.Printf(LogFormat, ERROR, l.ID, fmt.Sprintf(format, v...))
}

// Panicf print panic message
// then panic with input
func (l *Logger) Panicf(format string, v ...interface{}) {
	f := fmt.Sprintf(format, v...)
	log.Printf(LogFormat, PANIC, l.ID, f)
}

// Fatalf print fatal message
// then os.Exit(1)
func (l *Logger) Fatalf(format string, v ...interface{}) {
	log.Printf(LogFormat, FATAL, l.ID, fmt.Sprintf(format, v...))
	os.Exit(1)
}
