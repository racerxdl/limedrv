package limedrv

import (
	"fmt"
	"github.com/racerxdl/limedrv/limewrap"
)

// LMSDeviceAdvanced is a dummy structure just to separated the methods considered for "Advanced Usage"
// It does not have any data besides the methods to allow advanced settings of LMSDevice object.
type LMSDeviceAdvanced struct {
	parent *LMSDevice
}

// SetDigitalFilterTaps allows to manually set the GFIR digital filter taps from a channel.
// For enabling / disabling the GFIR when setting manual taps please use EnableGFIR / DisableGFIR in Advanced Section
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

// EnableGFIR enables a manually set GFIR Taps in the channel
func (d *LMSDeviceAdvanced) EnableGFir(gFirIdx, channelNumber int, isRX bool) {
	if limewrap.LMS_SetGFIR(d.parent.dev, !isRX, int64(channelNumber), limewrap.Lms_gfir_t(gFirIdx), true) != 0 {
		panic(fmt.Sprintf("Cannot enable GFir %s at %s: %s", d.parent.DeviceInfo.DeviceName, d.parent.DeviceInfo.Media, limewrap.LMS_GetLastErrorMessage()))
	}
}

// DisableGFIR disables a manually set GFIR Taps in the channel
func (d *LMSDeviceAdvanced) DisableGFir(gFirIdx, channelNumber int, isRX bool) {
	if limewrap.LMS_SetGFIR(d.parent.dev, !isRX, int64(channelNumber), limewrap.Lms_gfir_t(gFirIdx), false) != 0 {
		panic(fmt.Sprintf("Cannot disable GFir %s at %s: %s", d.parent.DeviceInfo.DeviceName, d.parent.DeviceInfo.Media, limewrap.LMS_GetLastErrorMessage()))
	}
}
