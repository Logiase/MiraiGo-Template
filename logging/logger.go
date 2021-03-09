package logging

// ILogger internal logger interface
//
// ILogger 内置Logger接口
// 内置默认实现 Logger
type ILogger interface {
	// Printf prints message
	Printf(format string, v ...interface{})

	// log functions
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Panicf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
}

// InternalLogger internal logger
var InternalLogger = NewLogger("internal")
