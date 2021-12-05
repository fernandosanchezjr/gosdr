package sdr

import (
	"github.com/fernandosanchezjr/gosdr/devices"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

var defaultScanTime = time.Second * 5

type Scanner struct {
	wg         sync.WaitGroup
	quitChan   chan struct{}
	DeviceChan chan []*devices.DeviceInfo
}

func NewScanner() *Scanner {
	var scanner = &Scanner{
		quitChan:   make(chan struct{}),
		DeviceChan: make(chan []*devices.DeviceInfo),
	}
	scanner.wg.Add(1)
	go scanner.loop()
	return scanner
}

func devicesDiffer(a []*devices.DeviceInfo, b []*devices.DeviceInfo) bool {
	if len(a) != len(b) {
		return true
	}
	for i := 0; i < len(a); i++ {
		if !a[i].Equals(b[i]) {
			return true
		}
	}
	return false
}

func (s *Scanner) loop() {
	var ticker = time.NewTicker(defaultScanTime)
	var knownDevices = ListDevices()
	log.WithField("count", len(knownDevices)).Println("Found devices")
	s.DeviceChan <- knownDevices
	for {
		select {
		case <-s.quitChan:
			s.wg.Done()
			return
		case <-ticker.C:
			var foundDevices = ListDevices()
			if devicesDiffer(knownDevices, foundDevices) {
				knownDevices = foundDevices
				log.WithField("count", len(foundDevices)).Println("Found devices")
				s.DeviceChan <- knownDevices
			}
		}
	}
}

func (s *Scanner) Close() {
	close(s.quitChan)
	s.wg.Wait()
}
