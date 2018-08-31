package limedrv

import (
	"fmt"
	"github.com/racerxdl/limedrv/limewrap"
	"strings"
	"unsafe"
)

type LMSDevice struct {
	dev uintptr
	initialized bool
	DeviceInfo DeviceInfo
	RXChannels []LMSChannel
	TXChannels []LMSChannel
	MinimumSampleRate float64
	MaximumSampleRate float64
}

func (d *LMSDevice) Close() {
	Close(d)
}

func (d *LMSDevice) init() {
	if limewrap.LMS_Reset(d.dev) != 0 {
		panic(fmt.Sprintf("Failed to reset %s at %s: %s", d.DeviceInfo.DeviceName, d.DeviceInfo.Media, limewrap.LMS_GetLastErrorMessage()))
	}

	if limewrap.LMS_Init(d.dev) != 0 {
		panic(fmt.Sprintf("Failed to init %s at %s: %s", d.DeviceInfo.DeviceName, d.DeviceInfo.Media, limewrap.LMS_GetLastErrorMessage()))
	}

	d.loadChannels()
}

func (d *LMSDevice) loadChannels() {
	// region Load RX Channels
	rxChannels := limewrap.LMS_GetNumChannels(d.dev, limewrap.LmsChRx)
	d.RXChannels = make([]LMSChannel, rxChannels)
	for i := 0; i < rxChannels; i++ {
		ch := LMSChannel{
			IsRX: true,
		}
		antennas := limewrap.LMS_GetAntennaList(d.dev, limewrap.LmsChRx, int64(i), nil)
		ch.Antennas = make([]LMSAntenna, antennas)

		if antennas > 0 {
			var nameArr [16*20]byte  // 16 bytes per lms_name_t
			var namePtr = (*string)(unsafe.Pointer(&nameArr[0]))
			limewrap.LMS_GetAntennaList(d.dev, limewrap.LmsChRx, int64(i), namePtr)
			for a := 0; a < antennas; a++ {
				var name = cleanString(string(nameArr[a*16:(a+1)*16]))
				var bw = createLms_range_t()
				limewrap.LMS_GetAntennaBW(d.dev, limewrap.LmsChRx, int64(i), int64(a), bw)
				ch.Antennas[a] = LMSAntenna{
					Name: name,
					Channel: i,
					MaximumFrequency: bw.GetMax(),
					MinimumFrequency: bw.GetMin(),
					Step: bw.GetStep(),
				}
			}
		}

		d.RXChannels[i] = ch
	}
	// endregion
	// region Load TX Channels
	txChannels := limewrap.LMS_GetNumChannels(d.dev, limewrap.LmsChTx)
	d.TXChannels = make([]LMSChannel, txChannels)
	for i := 0; i < txChannels; i++ {
		ch := LMSChannel{
			IsRX: false,
		}
		antennas := limewrap.LMS_GetAntennaList(d.dev, limewrap.LmsChTx, int64(i), nil)
		ch.Antennas = make([]LMSAntenna, antennas)

		if antennas > 0 {
			var nameArr [16*64]byte  // 16 bytes per lms_name_t
			var namePtr = (*string)(unsafe.Pointer(&nameArr[0]))
			limewrap.LMS_GetAntennaList(d.dev, limewrap.LmsChTx, int64(i), namePtr)
			for a := 0; a < antennas; a++ {
				var name = cleanString(string(nameArr[a*16:(a+1)*16]))
				var bw = createLms_range_t()
				limewrap.LMS_GetAntennaBW(d.dev, limewrap.LmsChTx, int64(i), int64(a), bw)
				ch.Antennas[a] = LMSAntenna{
					Name: name,
					Channel: i,
					MaximumFrequency: bw.GetMax(),
					MinimumFrequency: bw.GetMin(),
					Step: bw.GetStep(),
				}
			}
		}

		d.TXChannels[i] = ch
	}
	// endregion
}

func (d *LMSDevice) EnableChannel(channelNumber int, isRX bool) {
	if limewrap.LMS_EnableChannel(d.dev, ) != 0 {
		panic(fmt.Sprintf("Failed to enable channel in %s at %s: %s", d.DeviceInfo.DeviceName, d.DeviceInfo.Media, limewrap.LMS_GetLastErrorMessage()))
	}
}

func (d *LMSDevice) SetSampleRate(sampleRate float64, oversample int) {
	if limewrap.LMS_SetSampleRate(d.dev, sampleRate, int64(oversample)) != 0 {
		panic(fmt.Sprintf("Failed to set SampleRate to %f in %s at %s: %s", sampleRate, d.DeviceInfo.DeviceName, d.DeviceInfo.Media, limewrap.LMS_GetLastErrorMessage()))
	}
}

func (d *LMSDevice) GetSampleRate() (host float64, rf float64) {
	host = float64(0)
	rf = float64(0)
	//LMS_GetSampleRate (lms_device_t *device, bool dir_tx, size_t chan, float_type *host_Hz, float_type *rf_Hz)
	if limewrap.LMS_GetSampleRate(d.dev, limewrap.LmsChRx, 0, &host, &rf) != 0 {
		panic(fmt.Sprintf("Failed to get SampleRate in %s at %s: %s", d.DeviceInfo.DeviceName, d.DeviceInfo.Media, limewrap.LMS_GetLastErrorMessage()))
	}

	return host, rf
}

func (d *LMSDevice) String() string {
	var str = fmt.Sprintf("LMSDevice(%s)", d.DeviceInfo.DeviceName)

	str = fmt.Sprintf("%s\nRX Channels: %d", str, len(d.RXChannels))
	for i := 0; i < len(d.RXChannels); i++ {
		var chanStr = strings.Replace(d.RXChannels[i].String(), "\n", "\n\t", -1)
		str = fmt.Sprintf("%s\nChannel %d: %s", str, i, chanStr)
	}

	str = fmt.Sprintf("%s\nTX Channels: %d", str, len(d.TXChannels))
	for i := 0; i < len(d.TXChannels); i++ {
		var chanStr = strings.Replace(d.TXChannels[i].String(), "\n", "\n\t", -1)
		str = fmt.Sprintf("%s\nChannel %d: %s", str, i, chanStr)
	}

	return str
}