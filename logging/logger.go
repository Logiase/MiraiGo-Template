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
	LogFormat = "%s: %s"
)

// ILogger is internal logger interface
// You can implement by yourself or use internal logger
type ILogger interface {
	Printf(format string, v ...interface{})

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
}

func (_ *Logger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (_ *Logger) Debugf(format string, v ...interface{}) {
	log.Printf(LogFormat, DEBUG, fmt.Sprintf(format, v...))
}

func (_ *Logger) Infof(format string, v ...interface{}) {
	log.Printf(LogFormat, INFO, fmt.Sprintf(format, v...))
}

func (_ *Logger) Warnf(format string, v ...interface{}) {
	log.Printf(LogFormat, WARN, fmt.Sprintf(format, v...))
}

func (_ *Logger) Errorf(format string, v ...interface{}) {
	log.Printf(LogFormat, ERROR, fmt.Sprintf(format, v...))
}

func (_ *Logger) Panicf(format string, v ...interface{}) {
	f := fmt.Sprintf(format, v...)
	log.Printf(LogFormat, PANIC, f)
}

func (_ *Logger) Fatalf(format string, v ...interface{}) {
	log.Printf(LogFormat, FATAL, fmt.Sprintf(format, v...))
	os.Exit(1)
}
