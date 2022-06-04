package main

import (
	"github.com/fernandosanchezjr/gosdr/config"
	"github.com/fernandosanchezjr/gosdr/devices/sdr"
	"github.com/fernandosanchezjr/gosdr/utils"
	log "github.com/sirupsen/logrus"
)

func applicationLoop(manager *sdr.Manager) {
	go selector(manager)
}

func main() {
	config.ParseFlags()
	config.SetupLogger()
	log.Info("Starting GOSDRRX")
	var manager = sdr.NewManager()
	defer manager.Stop()
	go applicationLoop(manager)
	utils.Wait()
}
