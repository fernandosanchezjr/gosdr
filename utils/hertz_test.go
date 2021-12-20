package utils

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestParseHertz(t *testing.T) {
	var h, err = ParseHertz("24.9M")
	log.WithFields(log.Fields{
		"h":   h,
		"err": err,
	}).Info("Parsed")
}
