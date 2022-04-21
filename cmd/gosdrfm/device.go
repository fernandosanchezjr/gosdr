package main

import (
	"github.com/fernandosanchezjr/gosdr/devices"
	"github.com/fernandosanchezjr/gosdr/devices/sdr"
	log "github.com/sirupsen/logrus"
)

func deviceSelector(events chan sdr.DeviceEvent, deviceIds chan devices.Id) {
	for {
		select {
		case event := <-events:
			if event.EventType != sdr.DeviceAdded {
				continue
			}
			if !rtlSDR && event.Id.Type == devices.RTLSDR {
				continue
			}
			if deviceSerial != "" && deviceSerial == event.Id.Serial {
				log.WithFields(event.Id.Fields()).Trace("Selected device")
				deviceIds <- event.Id
				continue
			}
			if deviceIndex == event.Index {
				log.WithFields(event.Id.Fields()).WithField("index", deviceIndex).Trace(
					"Selected",
				)
				deviceIds <- event.Id
				continue
			}
		}
	}
}

func deviceController(manager *sdr.Manager, deviceIds chan devices.Id) {
	for {
		select {
		case id := <-deviceIds:
			var conn, connErr = manager.Open(id)
			if connErr != nil {
				log.WithFields(id.Fields()).WithError(connErr).Error("manager.Open")
				continue
			}
			if agcErr := conn.SetAGC(agc); agcErr != nil {
				log.WithFields(conn.Fields()).WithError(agcErr).Error("conn.SetAGC")
				manager.Close(id)
				continue
			}
			if autoGainErr := conn.SetAutoGain(autoGain); autoGainErr != nil {
				log.WithFields(conn.Fields()).WithError(autoGainErr).Error("conn.SetAutoGain")
				manager.Close(id)
				continue
			}
			if gainErr := conn.SetTunerGain(float32(gain)); gainErr != nil {
				log.WithFields(conn.Fields()).WithError(gainErr).Error("conn.SetTunerGain")
				manager.Close(id)
				continue
			}
			if ppmErr := conn.SetFrequencyCorrection(ppm); ppmErr != nil {
				log.WithFields(conn.Fields()).WithError(ppmErr).Error("conn.SetFrequencyCorrection")
				manager.Close(id)
				continue
			}
			if freqErr := conn.SetCenterFrequency(frequency); freqErr != nil {
				log.WithFields(conn.Fields()).WithError(freqErr).Error("conn.SetCenterFrequency")
				manager.Close(id)
				continue
			}
			if resetErr := conn.Reset(); resetErr != nil {
				log.WithFields(conn.Fields()).WithError(resetErr).Error("conn.Reset")
				manager.Close(id)
				continue
			}
			var bufferCount = conn.GetBuffersPerSecond()
			var input, graphErr = createGraph2(conn, bufferCount)
			if graphErr != nil {
				log.WithError(graphErr).Error("createGraph()")
				manager.Close(id)
				continue
			}
			manager.AddDeviceCleanup(id, func() {
				input.Close()
			})
			go sampleDevice(conn, bufferCount, input)
			log.WithFields(conn.Fields()).Info("Sampling SDR")
		}
	}
}
