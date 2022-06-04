package main

import (
	"github.com/fernandosanchezjr/gosdr/devices"
	"github.com/fernandosanchezjr/gosdr/devices/sdr"
	log "github.com/sirupsen/logrus"
)

func selector(manager *sdr.Manager) {
	for {
		select {
		case event := <-manager.DeviceChan:
			switch event.EventType {
			case sdr.DeviceRemoved:
				log.WithFields(event.Id.Fields()).WithField("index", event.Index).Debug(
					"Lost SDR",
				)
				continue
			case sdr.DeviceAdded:
				if !rtlSDR && event.Id.Type == devices.RTLSDR {
					continue
				}
				log.WithFields(event.Id.Fields()).WithField("index", event.Index).Debug(
					"Found SDR",
				)
				go controller(manager, event.Id)
			}
		}
	}
}

func controller(manager *sdr.Manager, id devices.Id) {
	var exitChan = make(chan struct{})
	var conn, connErr = manager.Open(id)
	if connErr != nil {
		log.WithFields(id.Fields()).WithError(connErr).Error("manager.Open")
		return
	}
	var logger = log.WithFields(conn.Fields())
	if resetErr := conn.Reset(); resetErr != nil {
		logger.WithError(resetErr).Error("conn.Reset")
		manager.Close(id)
		return
	}
	var bufferCount = conn.GetBuffersPerSecond() * bufferSeconds
	var input, graphErr = createGraph(conn, bufferCount)
	if graphErr != nil {
		logger.WithError(graphErr).Error("createGraph()")
		manager.Close(id)
		return
	}
	manager.AddDeviceCleanup(id, func() {
		input.Close()
		close(exitChan)
	})
	go sampleDevice(conn, bufferCount, input)
	logger.Info("SDR controller starting")
	for {
		select {
		case <-exitChan:
			logger.Info("SDR controller exiting")
			return
		}
	}
}
