package utils

import (
	"io/ioutil"
	"os"

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

func FileExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
