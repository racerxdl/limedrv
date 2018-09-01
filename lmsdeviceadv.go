package limedrv

import (
	"fmt"
	"github.com/racerxdl/limedrv/limewrap"
)

type LMSDeviceAdvanced struct {
	parent *LMSDevice
}

func (d *LMSDeviceAdvanced) SetDigitalFilterTaps(gFirIdx, channelNumber int, isRX bool, taps []float64) {
	if limewrap.LMS_SetGFIRCoeff(d.parent.dev, !isRX, int64(channelNumber), limewrap.Lms_gfir_t(gFirIdx), &taps[0], int64(len(taps))) != 0 {
		panic(fmt.Sprintf("Cannot set digital filter taps %s at %s: %s", d.parent.DeviceInfo.DeviceName, d.parent.DeviceInfo.Media, limewrap.LMS_GetLastErrorMessage()))
	}

	if isRX {
		d.parent.RXChannels[channelNumber].advancedFiltering = true
	} else {
		d.parent.TXChannels[channelNumber].advancedFiltering = true
	}
}

func (d *LMSDeviceAdvanced) EnableGFir(gFirIdx, channelNumber int, isRX bool) {
	if limewrap.LMS_SetGFIR(d.parent.dev, !isRX, int64(channelNumber), limewrap.Lms_gfir_t(gFirIdx), true) != 0 {
		panic(fmt.Sprintf("Cannot enable GFir %s at %s: %s", d.parent.DeviceInfo.DeviceName, d.parent.DeviceInfo.Media, limewrap.LMS_GetLastErrorMessage()))
	}
}

func (d *LMSDeviceAdvanced) DisableGFir(gFirIdx, channelNumber int, isRX bool) {
	if limewrap.LMS_SetGFIR(d.parent.dev, !isRX, int64(channelNumber), limewrap.Lms_gfir_t(gFirIdx), false) != 0 {
		panic(fmt.Sprintf("Cannot disable GFir %s at %s: %s", d.parent.DeviceInfo.DeviceName, d.parent.DeviceInfo.Media, limewrap.LMS_GetLastErrorMessage()))
	}
}
