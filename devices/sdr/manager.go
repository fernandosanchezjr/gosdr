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
	mtx           sync.Mutex
	wg            sync.WaitGroup
	quitChan      chan struct{}
	DeviceChan    chan DeviceEvent
	knownDevices  map[devices.Id]*devices.Info
	connections   map[devices.Id]devices.Connection
	deviceCleanup map[devices.Id][]func()
}

type DeviceEvent struct {
	EventType DeviceEventType
	Id        devices.Id
	Index     int
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
		"index":     de.Index,
	}
}

func NewManager() *Manager {
	var scanner = &Manager{
		quitChan:      make(chan struct{}),
		DeviceChan:    make(chan DeviceEvent, deviceEventBufferSize),
		knownDevices:  make(map[devices.Id]*devices.Info),
		connections:   make(map[devices.Id]devices.Connection),
		deviceCleanup: make(map[devices.Id][]func()),
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

func (s *Manager) processDevices(foundDevices []*devices.Info) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	var foundMap = toDeviceMap(foundDevices)
	var knownSet, foundSet = deviceIdSet(s.knownDevices), deviceIdSet(foundMap)
	var knownDiff = knownSet.Diff(foundSet)
	var foundDiff = foundSet.Diff(knownSet)
	for _, id := range knownDiff.Elements() {
		var deviceId = id.(devices.Id)
		var lostDeviceInfo = s.knownDevices[deviceId]
		go s.Close(deviceId)
		delete(s.knownDevices, deviceId)
		s.DeviceChan <- DeviceEvent{
			EventType: DeviceRemoved,
			Id:        deviceId,
			Index:     lostDeviceInfo.Index,
		}
		log.WithFields(lostDeviceInfo.Fields()).Trace("Lost device")
	}
	for _, id := range foundDiff.Elements() {
		var deviceId = id.(devices.Id)
		var foundDevice = foundMap[deviceId]
		s.knownDevices[deviceId] = foundDevice
		s.DeviceChan <- DeviceEvent{
			EventType: DeviceAdded,
			Id:        deviceId,
			Index:     foundDevice.Index,
		}
		log.WithFields(foundDevice.Fields()).Trace("Found device")
	}
}

func (s *Manager) loop() {
	var ticker = time.NewTicker(scanInterval)
	s.processDevices(ListDevices())
	for {
		select {
		case <-s.quitChan:
			s.wg.Done()
			return
		case <-ticker.C:
			s.processDevices(ListDevices())
		}
	}
}

func (s *Manager) Stop() {
	close(s.quitChan)
	s.wg.Wait()
	for id := range s.connections {
		s.Close(id)
	}
}

func (s *Manager) Open(id devices.Id) (devices.Connection, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if conn, found := s.connections[id]; found {
		return conn, nil
	}
	if info, found := s.knownDevices[id]; found {
		switch info.Type {
		case devices.RTLSDR:
			if conn, connErr := rtlsdr.OpenIndex(info.Index); connErr != nil {
				return nil, connErr
			} else {
				s.connections[id] = conn
				return conn, nil
			}
		default:
			return nil, fmt.Errorf("unknown device type: %d", info.Type)
		}
	} else {
		return nil, nil
	}
}

func (s *Manager) OpenAsync(id devices.Id) {
	go func() {
		if _, openErr := s.Open(id); openErr != nil {
			log.WithFields(id.Fields()).WithError(openErr).Error("OpenAsync")
		}
	}()
}

func (s *Manager) GetInfo(id devices.Id) (device *devices.Info, found bool) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	device, found = s.knownDevices[id]
	return
}

func (s *Manager) IsConnected(id devices.Id) bool {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	var _, found = s.connections[id]
	return found
}

func (s *Manager) Close(id devices.Id) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if conn, found := s.connections[id]; found {
		for _, cleanup := range s.deviceCleanup[id] {
			cleanup()
		}
		if closeErr := conn.Close(); closeErr != nil {
			log.WithFields(conn.Fields()).WithError(closeErr).Warn("Close")
		}
	}
	delete(s.connections, id)
}

func (s *Manager) CloseAsync(id devices.Id) {
	go s.Close(id)
}

func (s *Manager) AddDeviceCleanup(id devices.Id, f func()) {
	var cleanup = s.deviceCleanup[id]
	cleanup = append(cleanup, f)
	s.deviceCleanup[id] = cleanup
}
