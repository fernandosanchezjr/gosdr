package config

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestGetConfigDir(t *testing.T) {
	log.WithField("path", GetConfigDir()).Println("GetConfigDir")
}

func TestGetConfigFile(t *testing.T) {
	var handle, created, err = GetConfigFile("test.yaml")
	log.WithFields(log.Fields{
		"handle":  handle,
		"created": created,
		"err":     err,
	}).Println("GetConfigFile")
}
