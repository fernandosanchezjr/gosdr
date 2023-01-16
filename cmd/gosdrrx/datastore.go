package main

import (
	"github.com/fernandosanchezjr/gosdr/utils"
	log "github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
	"path"
)

func startDatastore(folder string) (*bolt.DB, error) {
	if folderErr := utils.CreateFolder(folder); folderErr != nil {
		return nil, folderErr
	}
	boltPath := path.Join(folder, "bolt")
	db, err := bolt.Open(boltPath, 0666, nil)
	if err == nil {
		log.WithField("path", boltPath).Info("Starting datastore")
	}
	return db, err
}

func stopDatastore(db *bolt.DB) {
	err := db.Close()
	if err != nil {
		log.WithError(err).Error("Error closing datastore")
	}
}

func listBuckets(db *bolt.DB) ([]string, error) {
	var buckets []string
	return buckets, db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			buckets = append(buckets, string(name))

			return nil
		})
	})
}
