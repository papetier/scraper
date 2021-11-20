package logger

import (
	"github.com/sirupsen/logrus"
	"log"
	"strings"
)

func Configure(level string) {
	var logrusLevel logrus.Level
	switch strings.ToLower(level) {
	case "panic":
		logrusLevel = logrus.PanicLevel
	case "fatal":
		logrusLevel = logrus.FatalLevel
	case "error":
		logrusLevel = logrus.ErrorLevel
	case "warn", "warning":
		logrusLevel = logrus.WarnLevel
	case "info":
		logrusLevel = logrus.InfoLevel
	case "debug":
		logrusLevel = logrus.DebugLevel
	case "trace":
		logrusLevel = logrus.TraceLevel
	default:
		log.Fatalf("incompatible log level: %s", level)
	}

	logrus.SetLevel(logrusLevel)
	logrus.Printf("log level set to: %s", logrusLevel)
}
