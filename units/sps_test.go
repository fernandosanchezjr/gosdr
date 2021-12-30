package units

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestParseSps(t *testing.T) {
	var s, err = ParseSps("24.9M")
	log.WithFields(log.Fields{
		"s":   s,
		"err": err,
	}).Info("Parsed")
}
