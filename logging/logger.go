package logging

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
)

var (
	DEBUG string = color.New(color.BgBlue).SprintfFunc()("DEBUG")
	INFO  string = color.New(color.BgGreen).SprintfFunc()("INFO ")
	WARN  string = color.New(color.BgYellow).SprintfFunc()("WARN ")
	ERROR string = color.New(color.BgRed).SprintfFunc()("ERROR")
	PANIC string = color.New(color.BgCyan).SprintfFunc()("PANIC")
	FATAL string = color.New(color.BgHiRed).SprintfFunc()("FATAL")
)

const (
	LogFormat = "%s ID(%s): %s"
)

// ILogger is internal logger interface
// You can implement by yourself or use internal logger
type ILogger interface {
	// Printf prints message
	Printf(format string, v ...interface{})

	// log interface
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Panicf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
}

// Logger is an internal implement of [ILogger]
// using standard "log"
type Logger struct {
	ID string
}

// NewLogger return a new Logger instance
func NewLogger(id string) Logger {
	return Logger{
		ID: id,
	}
}

func (l *Logger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	log.Printf(LogFormat, DEBUG, l.ID, fmt.Sprintf(format, v...))
}

func (l *Logger) Infof(format string, v ...interface{}) {
	log.Printf(LogFormat, INFO, l.ID, fmt.Sprintf(format, v...))
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	log.Printf(LogFormat, WARN, l.ID, fmt.Sprintf(format, v...))
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	log.Printf(LogFormat, ERROR, l.ID, fmt.Sprintf(format, v...))
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	f := fmt.Sprintf(format, v...)
	log.Printf(LogFormat, PANIC, l.ID, f)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	log.Printf(LogFormat, FATAL, l.ID, fmt.Sprintf(format, v...))
	os.Exit(1)
}
