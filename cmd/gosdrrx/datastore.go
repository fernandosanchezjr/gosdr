package main

import (
	"github.com/fernandosanchezjr/gosdr/config"
	"github.com/fernandosanchezjr/gosdr/utils"
	log "github.com/sirupsen/logrus"
	"github.com/xujiajun/nutsdb"
)

func startDatastore(config *config.Config) (*nutsdb.DB, error) {
	if folderErr := utils.CreateFolder(config.DataFolder); folderErr != nil {
		return nil, folderErr
	}
	var opt = nutsdb.DefaultOptions
	opt.Dir = config.DataFolder
	db, err := nutsdb.Open(opt)
	if err != nil {
		log.WithField("path", config.DataFolder).Info("Starting datastore")
	}
	return db, err
}

func stopDatastore(db *nutsdb.DB) {
	err := db.Close()
	log.WithError(err).Error("Error closing datastore")
}
