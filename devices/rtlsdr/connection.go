package rtlsdr

import (
	"github.com/fernandosanchezjr/gosdr/devices"
	"github.com/fernandosanchezjr/gosdr/units"
	rtl "github.com/jpoirier/gortlsdr"
	log "github.com/sirupsen/logrus"
)

const (
	defaultSampleRate = units.Sps(2400000)
	defaultBandwidth  = units.Hertz(2400000)
)

type Connection struct {
	Info            *devices.Info
	context         *rtl.Context
	Mode            rtl.SamplingMode
	PPM             int
	CenterFrequency units.Hertz
	OffsetTuning    bool
	SampleRate      units.Sps
	AGCMode         bool
	AutoGain        bool
	BiasTee         bool
	Gain            float32
	Gains           []float32
	TunerType       string
	RTLFrequency    units.Hertz
	TunerFrequency  units.Hertz
	TunerBandwidth  units.Hertz
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
	if !d.IsOpen() {
		return nil
	}
	var err = d.context.Close()
	if err != nil {
		log.WithError(err).WithFields(d.Info.Fields()).Error("context.Stop()")
	}
	d.context = nil
	return err
}

func (d *Connection) IsOpen() bool {
	return d.context != nil
}

func (d *Connection) init() error {
	if sampleErr := d.SetSampleRate(defaultSampleRate); sampleErr != nil {
		return sampleErr
	}
	return d.Refresh()
}

func (d *Connection) Refresh() (err error) {
	if d.Mode, err = d.context.GetDirectSampling(); err != nil {
		return
	}
	d.PPM = d.context.GetFreqCorrection()
	d.CenterFrequency = units.Hertz(d.context.GetCenterFreq())
	if d.OffsetTuning, err = d.context.GetOffsetTuning(); err != nil {
		return
	}
	d.SampleRate = units.Sps(d.context.GetSampleRate())
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
	d.RTLFrequency = units.Hertz(rtlFrequency)
	d.TunerFrequency = units.Hertz(tunerFrequency)
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
	f["autoGain"] = d.AutoGain
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

func (d *Connection) GetAGC() bool {
	return d.AGCMode
}

func (d *Connection) SetAGC(enabled bool) error {
	var err = d.context.SetAgcMode(enabled)
	if err == nil {
		d.AGCMode = enabled
	}
	return err
}

func (d *Connection) GetAutoGain() bool {
	return d.AutoGain
}

func (d *Connection) SetAutoGain(enabled bool) error {
	var err = d.context.SetTunerGainMode(!enabled)
	if err == nil {
		d.AutoGain = enabled
	}
	return err
}

func findNearestGain(gain float32, gains []float32) float32 {
	if len(gains) == 0 {
		return 0.0
	}
	var min, max = gains[0], gains[len(gains)-1]
	if gain < min {
		return min
	}
	if gain > max {
		return max
	}
	for i := 0; i < len(gains)-1; i++ {
		var current, next = gains[i], gains[i+1]
		if gain > next {
			continue
		}
		var diffLow, diffHigh = gain - current, next - gain
		if diffHigh <= diffLow {
			return next
		} else {
			return current
		}
	}
	return max
}

func (d *Connection) GetTunerGain() float32 {
	return d.Gain
}

func (d *Connection) SetTunerGain(gain float32) error {
	if d.AutoGain {
		return nil
	}
	var usedGain = findNearestGain(gain, d.Gains)
	var err = d.context.SetTunerGain(int(usedGain * 10.0))
	if err == nil {
		d.Gain = usedGain
	}
	return err
}

func (d *Connection) GetFrequencyCorrection() int {
	return d.PPM
}

func (d *Connection) SetFrequencyCorrection(ppm int) error {
	var err = d.context.SetFreqCorrection(ppm)
	if err != nil {
		d.PPM = ppm
	}
	return err
}

func (d *Connection) Reset() error {
	if resetBufferErr := d.context.ResetBuffer(); resetBufferErr != nil {
		return resetBufferErr
	}
	return nil
}

func (d *Connection) GetCenterFrequency() units.Hertz {
	return d.CenterFrequency
}

func (d *Connection) SetCenterFrequency(centerFrequency units.Hertz) error {
	var err = d.context.SetCenterFreq(int(centerFrequency))
	if err == nil {
		d.CenterFrequency = centerFrequency
	}
	return err
}

func (d *Connection) GetSampleRate() units.Sps {
	return d.SampleRate
}

func (d *Connection) SetSampleRate(sps units.Sps) error {
	var err = d.context.SetSampleRate(int(sps))
	if err == nil {
		d.SampleRate = sps
	}
	return err
}
