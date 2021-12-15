package sdr

import (
	"fmt"
	"github.com/fernandosanchezjr/gosdr/devices"
	"github.com/fernandosanchezjr/gosdr/devices/rtlsdr"
	log "github.com/sirupsen/logrus"
	"github.com/zoumo/goset"
	"sync"
	"time"
)

type DeviceEventType int

const (
	DeviceRemoved DeviceEventType = iota
	DeviceAdded
)

const scanInterval = time.Second * 1
const deviceEventBufferSize = 64

type Manager struct {
	wg           sync.WaitGroup
	quitChan     chan struct{}
	DeviceChan   chan DeviceEvent
	KnownDevices map[devices.Id]*devices.Info
}

type DeviceEvent struct {
	EventType DeviceEventType
	Id        devices.Id
}

func (det DeviceEventType) String() string {
	switch det {
	case DeviceRemoved:
		return "Device removed"
	case DeviceAdded:
		return "Device added"
	default:
		return "Unknown"
	}
}

func (de *DeviceEvent) Fields() log.Fields {
	return log.Fields{
		"eventType": de.EventType,
		"type":      de.Id.Type,
		"serial":    de.Id.Serial,
	}
}

func NewManager() *Manager {
	var scanner = &Manager{
		quitChan:     make(chan struct{}),
		DeviceChan:   make(chan DeviceEvent, deviceEventBufferSize),
		KnownDevices: make(map[devices.Id]*devices.Info),
	}
	scanner.wg.Add(1)
	go scanner.loop()
	return scanner
}

func toDeviceMap(deviceList []*devices.Info) (deviceMap map[devices.Id]*devices.Info) {
	deviceMap = make(map[devices.Id]*devices.Info)
	for _, device := range deviceList {
		deviceMap[device.Id()] = device
	}
	return
}

func deviceIdSet(deviceMap map[devices.Id]*devices.Info) (ids goset.Set) {
	ids = goset.NewSet()
	for deviceId := range deviceMap {
		if err := ids.Add(deviceId); err != nil {
			log.WithError(err).Warn("deviceIdSet")
		}
	}
	return
}

func (s *Manager) calculateDifferences(foundDevices []*devices.Info) {
	var foundMap = toDeviceMap(foundDevices)
	var knownSet, foundSet = deviceIdSet(s.KnownDevices), deviceIdSet(foundMap)
	var knownDiff = knownSet.Diff(foundSet)
	var foundDiff = foundSet.Diff(knownSet)
	for _, id := range knownDiff.Elements() {
		var deviceId = id.(devices.Id)
		var lostDeviceInfo = s.KnownDevices[deviceId]
		delete(s.KnownDevices, deviceId)
		s.DeviceChan <- DeviceEvent{
			EventType: DeviceRemoved,
			Id:        deviceId,
		}
		log.WithFields(lostDeviceInfo.Fields()).Println("Lost device")
	}
	for _, id := range foundDiff.Elements() {
		var deviceId = id.(devices.Id)
		var foundDevice = foundMap[deviceId]
		s.KnownDevices[deviceId] = foundDevice
		s.DeviceChan <- DeviceEvent{
			EventType: DeviceAdded,
			Id:        deviceId,
		}
		log.WithFields(foundDevice.Fields()).Println("Found device")
	}
}

func (s *Manager) loop() {
	var ticker = time.NewTicker(scanInterval)
	s.calculateDifferences(ListDevices())
	for {
		select {
		case <-s.quitChan:
			s.wg.Done()
			return
		case <-ticker.C:
			s.calculateDifferences(ListDevices())
		}
	}
}

func (s *Manager) Close() {
	close(s.quitChan)
	s.wg.Wait()
}

func (s *Manager) Open(info *devices.Info) (devices.Connection, error) {
	switch info.Type {
	case devices.RTLSDR:
		return rtlsdr.OpenIndex(info.Index)
	default:
		return nil, fmt.Errorf("unknown device type: %d", info.Type)
	}
}
