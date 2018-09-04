package limedrv

import (
	"fmt"
)

// LMSAntenna is a struct that represents the Antenna Port information
type LMSAntenna struct {
	Name             string
	Channel          int
	MinimumFrequency float64
	MaximumFrequency float64
	Step             float64

	parent *LMSChannel
	index  int
}

// Set sets this antenna port as the default in parent channel
func (a *LMSAntenna) Set() {
	a.parent.parent.SetAntenna(a.index, a.parent.parentIndex, a.parent.IsRX)
}

// String returns a representation of the antenna port data
func (a *LMSAntenna) String() string {
	return fmt.Sprintf("%6s: %14.0f -> %14.0f Hz", a.Name, a.MinimumFrequency, a.MaximumFrequency)
}
