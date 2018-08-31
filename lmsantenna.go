package limedrv

import "fmt"

type LMSChannel struct {
	Antennas []LMSAntenna
	IsRX bool
}

type LMSAntenna struct {
	Name string
	Channel int
	MinimumFrequency float64
	MaximumFrequency float64
	Step float64
}

func (a *LMSAntenna) String() string {
	return fmt.Sprintf("%4s: %14.0f -> %14.0f Hz", a.Name, a.MinimumFrequency, a.MaximumFrequency)
}

func (c *LMSChannel) String() string {
	var str = fmt.Sprintf("\nIs RX: %t\nAntennas: %d", c.IsRX, len(c.Antennas))
	for i := 0; i < len(c.Antennas); i++ {
		str = fmt.Sprintf("%s\n\t%s", str, c.Antennas[i].String())
	}

	return str
}