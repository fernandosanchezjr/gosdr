package units

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestParseSps(t *testing.T) {
	var s, err = ParseSps("2.4M")
	log.WithFields(log.Fields{
		"s":   s,
		"err": err,
	}).Info("Parsed")
}

func TestSps_NearestSize(t *testing.T) {
	var s = Sps(2400000)
	if s.NearestSize(512) != Sps(2400256) {
		t.Fail()
	}
}
