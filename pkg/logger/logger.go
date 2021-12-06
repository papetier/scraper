package logger

import (
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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

	// save logs to files
	setupLogFiles()
}

func setupLogFiles() {
	var pathMap lfshook.PathMap
	pathMap = make(map[logrus.Level]string)

	if viper.GetString("LOG_TRACE_FILE") != "" {
		pathMap[logrus.TraceLevel] = viper.GetString("LOG_TRACE_FILE")
	}
	if viper.GetString("LOG_DEBUG_FILE") != "" {
		pathMap[logrus.DebugLevel] = viper.GetString("LOG_DEBUG_FILE")
	}
	if viper.GetString("LOG_INFO_FILE") != "" {
		pathMap[logrus.InfoLevel] = viper.GetString("LOG_INFO_FILE")
	}
	if viper.GetString("LOG_WARN_FILE") != "" {
		pathMap[logrus.WarnLevel] = viper.GetString("LOG_WARN_FILE")
	}
	if viper.GetString("LOG_ERROR_FILE") != "" {
		pathMap[logrus.ErrorLevel] = viper.GetString("LOG_ERROR_FILE")
	}
	if viper.GetString("LOG_FATAL_FILE") != "" {
		pathMap[logrus.FatalLevel] = viper.GetString("LOG_FATAL_FILE")
	}
	if viper.GetString("LOG_PANIC_FILE") != "" {
		pathMap[logrus.PanicLevel] = viper.GetString("LOG_PANIC_FILE")
	}

	if len(pathMap) > 0 {
		logrus.StandardLogger().Hooks.Add(lfshook.NewHook(
			pathMap,
			&logrus.JSONFormatter{},
		))
	}
}
