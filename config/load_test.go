package config

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestGetConfigDir(t *testing.T) {
	log.WithField("path", GetConfigDir()).Println("GetConfigDir")
}

func TestLoadConfigFile(t *testing.T) {
	var handle, created, err = LoadConfigFile("test.yaml")
	log.WithFields(log.Fields{
		"handle":  handle,
		"created": created,
		"err":     err,
	}).Println("LoadConfigFile")
}

func TestLoadConfig(t *testing.T) {
	var config, loadErr = LoadConfig("test")
	log.WithFields(log.Fields{
		"config": config,
		"err":    loadErr,
	}).Println("LoadConfig")
}
