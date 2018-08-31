package limedrv

import (
	"github.com/racerxdl/limedrv/limewrap"
	"strings"
	"unsafe"
)

func cleanString(s string) string {
	return strings.Trim(s, "\u0000 ")
}

type lranget struct {
	min float64
	max float64
	step float64
}

func (l *lranget) Swigcptr() uintptr {
	return (uintptr)(unsafe.Pointer(l))
}

func (l *lranget) SwigIsLms_range_t() {}
func (l *lranget) SetMin(min float64) { l.min = min }
func (l *lranget) GetMin() float64 { return l.min }
func (l *lranget) SetMax(max float64) { l.max = max }
func (l *lranget) GetMax() float64 { return l.max }
func (l *lranget) SetStep(step float64) { l.step = step }
func (l *lranget) GetStep() float64 { return l.step }


func createLms_range_t() limewrap.Lms_range_t {
	return &lranget{}
}

func idev2dev(deviceinfo i_deviceinfo) DeviceInfo {
	var deviceStr = string(deviceinfo.DeviceName[:64])
	var z = strings.Split(deviceStr, ",")

	var DeviceName string
	var Media string
	var Module string
	var Addr string
	var Serial string

	for i := 0; i < len(z); i++ {
		var k = strings.Split(z[i], "=")
		if len(k) == 1 {
			DeviceName = k[0]
		} else {
			switch strings.ToLower(strings.Trim(k[0], " ")) {
			case "media": Media = cleanString(k[1]); break
			case "module": Module = cleanString(k[1]); break
			case "addr": Addr = cleanString(k[1]); break
			case "serial": Serial = cleanString(k[1]); break
			}
		}
	}


	return DeviceInfo{
		DeviceName: DeviceName,
		Media: Media,
		Module: Module,
		Addr: Addr,
		Serial: Serial,
		FirmwareVersion: cleanString(string(deviceinfo.FirmwareVersion[:16])),
		HardwareVersion: cleanString(string(deviceinfo.HardwareVersion[:16])),
		GatewareVersion: cleanString(string(deviceinfo.GatewareVersion[:16])),
		GatewareTargetBoard: cleanString(string(deviceinfo.GatewareTargetBoard[:16])),
		origDevInfo: deviceinfo,
	}
}
