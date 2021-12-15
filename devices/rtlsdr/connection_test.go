package rtlsdr

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestOpen(t *testing.T) {
	var device, openErr = OpenIndex(0)
	if openErr != nil {
		t.Fatal(openErr)
	}
	log.WithFields(device.Fields()).Println("Opened")
	var closeErr = device.Close()
	if closeErr != nil {
		t.Fatal(closeErr)
	}
}
