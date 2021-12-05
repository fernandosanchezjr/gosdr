package rtlsdr

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestListDevices(t *testing.T) {
	var devices = ListDevices()
	for _, deviceInfo := range devices {
		log.WithFields(deviceInfo.Fields()).Println("ListDevices")
	}
}
