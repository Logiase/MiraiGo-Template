package utils

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

func ReadFile(path string) []byte {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.WithError(err).WithField("util", "ReadFile").Errorf("unable to read '%s'", path)
		return nil
	}
	return bytes
}
