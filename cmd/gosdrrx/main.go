package main

import (
	"github.com/fernandosanchezjr/gosdr/config"
	"github.com/fernandosanchezjr/gosdr/devices/sdr"
	"github.com/fernandosanchezjr/gosdr/utils"
	log "github.com/sirupsen/logrus"
)

const appName = "gosdrrx"

func main() {
	log.Info("Starting ", appName)
	config.ParseFlags()
	config.SetupLogger()
	var cfg, configErr = config.LoadConfig(appName)
	if configErr != nil {
		log.WithError(configErr).Fatal("Error reading configuration")
	}
	var db, dbErr = startDatastore(cfg)
	if dbErr != nil {
		log.WithError(dbErr).Fatal("Error opening datastore")
	}
	defer stopDatastore(db)
	var manager = sdr.NewManager()
	defer manager.Stop()
	go selector(manager)
	utils.Wait()
}
