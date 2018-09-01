package limedrv

import "github.com/racerxdl/limedrv/limewrap"

type LMSChannel struct {
	Antennas                []LMSAntenna
	IsRX                    bool
	parent                  *LMSDevice
	parentIndex             int
	stream                  limewrap.Lms_stream_t
	currentDigitalBandwidth float64
	digitalFilterEnabled    bool
	advancedFiltering       bool
}

func (c *LMSChannel) Enable() *LMSChannel {
	c.parent.EnableChannel(c.parentIndex, c.IsRX)
	return c
}

func (c *LMSChannel) Disable() *LMSChannel {
	c.parent.DisableChannel(c.parentIndex, c.IsRX)
	return c
}

func (c *LMSChannel) SetGainDB(gain uint) *LMSChannel {
	c.parent.SetGainDB(c.parentIndex, c.IsRX, gain)
	return c
}

func (c *LMSChannel) SetGainNormalized(gain float64) *LMSChannel {
	c.parent.SetGainNormalized(c.parentIndex, c.IsRX, gain)
	return c
}

func (c *LMSChannel) GetGainDB() uint {
	return c.parent.GetGainDB(c.parentIndex, c.IsRX)
}

func (c *LMSChannel) GetGainNormalized() float64 {
	return c.parent.GetGainNormalized(c.parentIndex, c.IsRX)
}

func (c *LMSChannel) SetLPF(bandwidth float64) *LMSChannel {
	c.parent.SetLPF(c.parentIndex, c.IsRX, bandwidth)
	return c
}

func (c *LMSChannel) GetLPF() float64 {
	return c.parent.GetLPF(c.parentIndex, c.IsRX)
}

func (c *LMSChannel) EnableLPF() *LMSChannel {
	c.parent.EnableLPF(c.parentIndex, c.IsRX)
	return c
}

func (c *LMSChannel) DisableLPF() *LMSChannel {
	c.parent.EnableLPF(c.parentIndex, c.IsRX)
	return c
}

func (c *LMSChannel) SetAntenna(idx int) *LMSChannel {
	c.parent.SetAntenna(idx, c.parentIndex, c.IsRX)
	return c
}

func (c *LMSChannel) SetAntennaByName(name string) *LMSChannel {
	c.parent.SetAntennaByName(name, c.parentIndex, c.IsRX)
	return c
}

func (c *LMSChannel) start() {
	if c.stream != nil {
		limewrap.LMS_StartStream(c.stream)
	}
}

func (c *LMSChannel) stop() {
	if c.stream != nil {
		limewrap.LMS_StopStream(c.stream)
	}
}

func (c *LMSChannel) SetCenterFrequency(centerFrequency float64) *LMSChannel {
	c.parent.SetCenterFrequency(c.parentIndex, c.IsRX, centerFrequency)
	return c
}

func (c *LMSChannel) GetCenterFrequency() float64 {
	return c.parent.GetCenterFrequency(c.parentIndex, c.IsRX)
}
