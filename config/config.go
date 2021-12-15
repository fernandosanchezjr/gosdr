package config

import (
	"github.com/kirsle/configdir"
	"os"
	"path/filepath"
)

func GetConfigDir() string {
	return configdir.LocalConfig("gosdr/v0")
}

func GetConfigFile(name string) (handle *os.File, created bool, createErr error) {
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
