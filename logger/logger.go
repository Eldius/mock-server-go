package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Entry

func init() {
	hostname, _ := os.Hostname()
	var standardFields = logrus.Fields{
		"hostname": hostname,
		"app":      "mock-server",
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetReportCaller(true)
	log = logrus.WithFields(standardFields)
}

func Log() *logrus.Entry {
	return log
}
