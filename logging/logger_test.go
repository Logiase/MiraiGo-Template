package logging_test

import (
	"github.com/Logiase/MiraiGo-Template/v2/logging"
	"testing"
)

func TestLogger(t *testing.T) {
	logger := logging.NewLogger("test")

	logger.Printf("print\n")
	logger.Debugf("debug\n")
	logger.Infof("info\n")
	logger.Warnf("warn\n")
	logger.Errorf("error\n")
	defer func() {
		if recover(); true {
			logger.Fatalf("fatal\n")
		}
	}()
	logger.Panicf("panic\n")
}
