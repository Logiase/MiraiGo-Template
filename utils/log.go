package utils

import (
	"github.com/sirupsen/logrus"
)

// NewModuleLogger - 提供一个为 Module 使用的 logrus.Entry
// 包含 logrus.Fields
func NewModuleLogger(name string) *logrus.Entry {
	return logrus.WithField("module", name)
}
