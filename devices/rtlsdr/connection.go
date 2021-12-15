package rtlsdr

import (
	"github.com/fernandosanchezjr/gosdr/devices"
	"github.com/fernandosanchezjr/gosdr/utils"
	rtl "github.com/jpoirier/gortlsdr"
	log "github.com/sirupsen/logrus"
)

const (
	defaultSampleRate = utils.Hertz(2400000)
	defaultBandwidth  = utils.Hertz(2400000)
)

type Connection struct {
	Info            *devices.Info
	context         *rtl.Context
	Mode            rtl.SamplingMode
	PPM             int
	CenterFrequency utils.Hertz
	OffsetTuning    bool
	SampleRate      utils.Hertz
	AGCMode         bool
	BiasTee         bool
	Gain            float32
	Gains           []float32
	TunerType       string
	RTLFrequency    utils.Hertz
	TunerFrequency  utils.Hertz
	TunerBandwidth  utils.Hertz
}

func OpenIndex(index int) (*Connection, error) {
	var context, openErr = rtl.Open(index)
	if openErr != nil {
		log.WithError(openErr).WithField("index", index).Error("rtl.IsOpen")
		return nil, openErr
	}
	var info, infoErr = GetInfo(index)
	if infoErr != nil {
		return nil, infoErr
	}
	var device = &Connection{
		Info:           info,
		context:        context,
		SampleRate:     defaultSampleRate,
		TunerBandwidth: defaultBandwidth,
	}
	if infoErr := device.init(); infoErr != nil {
		return nil, infoErr
	}
	return device, nil
}

func (d *Connection) Close() error {
	var err = d.context.Close()
	if err != nil {
		log.WithError(err).WithFields(d.Info.Fields()).Error("context.Close()")
	}
	d.context = nil
	return err
}

func (d *Connection) IsOpen() bool {
	return d.context != nil
}

func (d *Connection) init() (err error) {
	return d.Refresh()
}

func (d *Connection) Refresh() (err error) {
	if d.Mode, err = d.context.GetDirectSampling(); err != nil {
		return
	}
	d.PPM = d.context.GetFreqCorrection()
	d.CenterFrequency = utils.Hertz(d.context.GetCenterFreq())
	if d.OffsetTuning, err = d.context.GetOffsetTuning(); err != nil {
		return
	}
	d.SampleRate = utils.Hertz(d.context.GetSampleRate())
	d.Gain = float32(d.context.GetTunerGain()) / 10.0
	var gains []int
	if gains, err = d.context.GetTunerGains(); err != nil {
		return err
	}
	d.Gains = make([]float32, len(gains))
	for pos, g := range gains {
		d.Gains[pos] = float32(g) / 10.0
	}
	d.TunerType = d.context.GetTunerType()
	var rtlFrequency, tunerFrequency int
	if rtlFrequency, tunerFrequency, err = d.context.GetXtalFreq(); err != nil {
		return
	}
	d.RTLFrequency = utils.Hertz(rtlFrequency)
	d.TunerFrequency = utils.Hertz(tunerFrequency)
	return
}

func (d *Connection) Fields() log.Fields {
	var f = d.Info.Fields()
	f["open"] = d.IsOpen()
	f["mode"] = d.Mode
	f["ppm"] = d.PPM
	f["centerFrequency"] = d.CenterFrequency
	f["offsetTuning"] = d.OffsetTuning
	f["sampleRate"] = d.SampleRate
	f["agcMode"] = d.AGCMode
	f["biasTee"] = d.BiasTee
	f["gain"] = d.Gain
	f["gains"] = d.Gains
	f["tunerType"] = d.TunerType
	f["rtlFrequency"] = d.RTLFrequency
	f["tunerFrequency"] = d.TunerFrequency
	return f
}

func (d *Connection) GetInfo() *devices.Info {
	return d.Info
}
