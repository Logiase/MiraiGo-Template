package logging_test

import (
	"testing"

	"github.com/Logiase/MiraiGo-Template/v2/logging"
)

func TestLogger(t *testing.T) {
	logger := logging.NewLogger("test")

	logger.Printf("print")
	logger.Debugf("debug")
	logger.Infof("info")
	logger.Warnf("warn")
	logger.Errorf("error")
	defer func() {
		if recover(); true {
			logger.Fatalf("fatal")
		}
	}()
	// logger.Panicf("panic")
}
