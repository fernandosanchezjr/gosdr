package config

import (
	"flag"
	"github.com/kirsle/configdir"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
)

var logFile *os.File
var logPath string
var logToFile = true
var logLevel = logrus.InfoLevel
var logLevelStr string

func init() {
	logPath = path.Join(GetConfigDir(), "logs", "log.out")
	flag.BoolVar(&logToFile, "log-to-file", logToFile, "enable logging to file")
	flag.StringVar(&logLevelStr, "log-level", logLevel.String(),
		"log level (panic, fatal, error, warn, info, debug, trace)")
	flag.StringVar(&logPath, "log-path", logPath, "log path")
}

func getLogFile() (*os.File, error) {
	logFolder := path.Base(logPath)
	if createErr := configdir.MakePath(logFolder); createErr != nil {
		return nil, createErr
	}
	return os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
}

func exitHandler() {
	if logFile != nil {
		_ = logFile.Close()
	}
}

func SetupLogger() {
	var err error
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true, FullTimestamp: true})
	logrus.RegisterExitHandler(exitHandler)
	if logLevel, err = logrus.ParseLevel(logLevelStr); err != nil {
		logrus.WithField("log-level", logLevelStr).Error("Invalid")
		logLevel = logrus.InfoLevel
	}
	logrus.SetLevel(logLevel)
	var details = logrus.Fields{
		"log-level":   logLevel,
		"log-to-file": logToFile,
	}
	if logToFile {
		if logFile, err = getLogFile(); err != nil {
			logrus.WithError(err).Error("getLogFile")
		} else {
			logrus.SetOutput(io.MultiWriter(logFile, os.Stdout))
			details["output"] = logPath
		}
	}
	logrus.WithFields(details).Info("Logging")
}
