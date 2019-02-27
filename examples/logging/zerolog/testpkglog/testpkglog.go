package testpkglog

import (
	log "github.com/sirupsen/logrus"
)

var logger *log.Entry

func init() {
	logger = log.WithFields(log.Fields{"pkg": "tesgpkglog"})
}

func Foo() {
	logger.Debug("Here we are in the package")
}
