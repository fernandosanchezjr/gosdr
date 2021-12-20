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

func init() {
	flag.BoolVar(&logToFile, "log-to-file", logToFile, "enable logging to file")
}

func getLogFile() (*os.File, error) {
	logFolder := path.Join(GetConfigDir(), "logs")
	if createErr := configdir.MakePath(logFolder); createErr != nil {
		return nil, createErr
	}
	logPath = path.Join(logFolder, "log.out")
	return os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
}

func exitHandler() {
	if logFile != nil {
		_ = logFile.Close()
	}
}

func SetupLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true, FullTimestamp: true})
	logrus.RegisterExitHandler(exitHandler)
	logrus.SetLevel(logrus.DebugLevel)
	if logToFile {
		var err error
		if logFile, err = getLogFile(); err != nil {
			logrus.WithError(err).Error("getLogFile")
		} else {
			logrus.SetOutput(io.MultiWriter(logFile, os.Stdout))
			logrus.WithField("output", logPath).Info("Logging")
		}
	}
}
