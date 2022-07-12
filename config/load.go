package config

import (
	"github.com/kirsle/configdir"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func GetConfigDir() string {
	return configdir.LocalConfig("gosdr/v0")
}

func LoadConfigFile(name string) (handle *os.File, created bool, createErr error) {
	var configDir = GetConfigDir()
	if createErr = configdir.MakePath(configDir); createErr != nil {
		return
	}
	var configPath = filepath.Join(configDir, name)
	if _, statErr := os.Stat(configPath); os.IsNotExist(statErr) {
		handle, createErr = os.OpenFile(configPath, os.O_CREATE|os.O_RDWR, 0600)
		created = true
	} else {
		handle, createErr = os.OpenFile(configPath, os.O_RDWR, 0600)
	}
	return
}

func LoadConfig(name string) (*Config, error) {
	var configFile, _, loadErr = LoadConfigFile(strings.Join([]string{name, "yml"}, "."))
	if loadErr != nil {
		return nil, loadErr
	}
	log.WithFields(log.Fields{
		"path": configFile.Name(),
	}).Println("Loading config")
	var data, readErr = ioutil.ReadAll(configFile)
	if readErr != nil {
		return nil, readErr
	}
	c := &Config{}
	if unmarshalErr := yaml.Unmarshal(data, c); unmarshalErr != nil {
		return nil, unmarshalErr
	}
	return c, nil
}
