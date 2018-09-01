package limedrv

import (
	"fmt"
	"github.com/racerxdl/limedrv/limewrap"
)

type LMSChannel struct {
	Antennas    []LMSAntenna
	IsRX        bool
	parent      *LMSDevice
	parentIndex int
	stream 		limewrap.Lms_stream_t
}

type LMSAntenna struct {
	Name             string
	Channel          int
	MinimumFrequency float64
	MaximumFrequency float64
	Step             float64
	parent           *LMSChannel
	index            int
}

func (c *LMSChannel) Enable() {
	c.parent.EnableChannel(c.parentIndex, c.IsRX)
}

func (c *LMSChannel) Disable() {
	c.parent.DisableChannel(c.parentIndex, c.IsRX)
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

func (a *LMSAntenna) Set() {
	a.parent.parent.SetAntenna(a.index, a.parent.parentIndex, a.parent.IsRX)
}

func (a *LMSAntenna) String() string {
	return fmt.Sprintf("%6s: %14.0f -> %14.0f Hz", a.Name, a.MinimumFrequency, a.MaximumFrequency)
}

func (c *LMSChannel) String() string {
	var str = fmt.Sprintf("\nIs RX: %t\nAntennas: %d", c.IsRX, len(c.Antennas))
	for i := 0; i < len(c.Antennas); i++ {
		str = fmt.Sprintf("%s\n\t%s", str, c.Antennas[i].String())
	}

	return str
}
