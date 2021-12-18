package utils

import "github.com/sirupsen/logrus"

// Close safety closes with log
func Close(closeFn func() error) {
	if err := closeFn(); err != nil {
		logrus.Error(err)
	}
}
