package limedrv

import (
	"fmt"
	"github.com/racerxdl/limedrv/limewrap"
	"os"
	"strings"
	"unsafe"
)

type LMSDevice struct {
	dev               uintptr
	initialized       bool
	DeviceInfo        DeviceInfo
	RXChannels        []*LMSChannel
	TXChannels        []*LMSChannel
	MinimumSampleRate float64
	MaximumSampleRate float64
	IQFormat          int

	RXLPFMaxFrequency float64
	RXLPFMinFrequency float64
	TXLPFMaxFrequency float64
	TXLPFMinFrequency float64

	Advanced LMSDeviceAdvanced

	controlChan chan bool
	running     bool
	callback    func([]complex64, int, uint64)
}

// region Private Methods

func (d *LMSDevice) init() {
	if limewrap.LMS_Reset(d.dev) != 0 {
		panic(fmt.Sprintf("Failed to reset %s at %s: %s", d.DeviceInfo.DeviceName, d.DeviceInfo.Media, limewrap.LMS_GetLastErrorMessage()))
	}

	if limewrap.LMS_Init(d.dev) != 0 {
		panic(fmt.Sprintf("Failed to init %s at %s: %s", d.DeviceInfo.DeviceName, d.DeviceInfo.Media, limewrap.LMS_GetLastErrorMessage()))
	}

	d.loadChannels()
	d.initSampleRateRange()
}

func (d *LMSDevice) loadChannels() {
	// region Load RX Channels
	var bw = createLms_range_t()
	limewrap.LMS_GetLPFBWRange(d.dev, limewrap.LmsChRx, bw)

	d.RXLPFMaxFrequency = bw.GetMax()
	d.RXLPFMinFrequency = bw.GetMin()

	rxChannels := limewrap.LMS_GetNumChannels(d.dev, limewrap.LmsChRx)
	d.RXChannels = make([]*LMSChannel, rxChannels)
	for i := 0; i < rxChannels; i++ {
		ch := LMSChannel{
			IsRX:              true,
			parent:            d,
			parentIndex:       i,
			advancedFiltering: false,
		}
		antennas := limewrap.LMS_GetAntennaList(d.dev, limewrap.LmsChRx, int64(i), nil)
		ch.Antennas = make([]LMSAntenna, antennas)

		if antennas > 0 {
			var nameArr [16 * 20]byte // 16 bytes per lms_name_t
			var namePtr = (*string)(unsafe.Pointer(&nameArr[0]))
			limewrap.LMS_GetAntennaList(d.dev, limewrap.LmsChRx, int64(i), namePtr)
			for a := 0; a < antennas; a++ {
				var name = cleanString(string(nameArr[a*16 : (a+1)*16]))
				limewrap.LMS_GetAntennaBW(d.dev, limewrap.LmsChRx, int64(i), int64(a), bw)
				ch.Antennas[a] = LMSAntenna{
					Name:             name,
					Channel:          i,
					MaximumFrequency: bw.GetMax(),
					MinimumFrequency: bw.GetMin(),
					Step:             bw.GetStep(),
					parent:           &ch,
					index:            a,
				}
			}
		}

		d.RXChannels[i] = &ch
	}
	// endregion
	// region Load TX Channels
	limewrap.LMS_GetLPFBWRange(d.dev, limewrap.LmsChTx, bw)

	d.TXLPFMaxFrequency = bw.GetMax()
	d.TXLPFMinFrequency = bw.GetMin()

	txChannels := limewrap.LMS_GetNumChannels(d.dev, limewrap.LmsChTx)
	d.TXChannels = make([]*LMSChannel, txChannels)
	for i := 0; i < txChannels; i++ {
		ch := LMSChannel{
			IsRX:              false,
			parent:            d,
			parentIndex:       i,
			advancedFiltering: false,
		}
		antennas := limewrap.LMS_GetAntennaList(d.dev, limewrap.LmsChTx, int64(i), nil)
		ch.Antennas = make([]LMSAntenna, antennas)

		if antennas > 0 {
			var nameArr [16 * 64]byte // 16 bytes per lms_name_t
			var namePtr = (*string)(unsafe.Pointer(&nameArr[0]))
			limewrap.LMS_GetAntennaList(d.dev, limewrap.LmsChTx, int64(i), namePtr)
			for a := 0; a < antennas; a++ {
				var name = cleanString(string(nameArr[a*16 : (a+1)*16]))
				limewrap.LMS_GetAntennaBW(d.dev, limewrap.LmsChTx, int64(i), int64(a), bw)
				ch.Antennas[a] = LMSAntenna{
					Name:             name,
					Channel:          i,
					MaximumFrequency: bw.GetMax(),
					MinimumFrequency: bw.GetMin(),
					Step:             bw.GetStep(),
					parent:           &ch,
					index:            a,
				}
			}
		}

		d.TXChannels[i] = &ch
	}
	// endregion
}

func (d *LMSDevice) initSampleRateRange() {
	var bw = createLms_range_t()
	if limewrap.LMS_GetSampleRateRange(d.dev, limewrap.LmsChRx, bw) != 0 {
		panic(fmt.Sprintf("Failed to get sample rate range %s at %s: %s", d.DeviceInfo.DeviceName, d.DeviceInfo.Media, limewrap.LMS_GetLastErrorMessage()))
	}

	d.MinimumSampleRate = bw.GetMin()
	d.MaximumSampleRate = bw.GetMax()
	d.SetSampleRate(1e6, 4)
}

func (d *LMSDevice) setupStream(channelNumber int, isRX bool) {
	var ch *LMSChannel

	if isRX {
		ch = d.RXChannels[channelNumber]
	} else {
		ch = d.TXChannels[channelNumber]
	}

	if ch.stream != nil {
		limewrap.LMS_DestroyStream(d.dev, ch.stream)
		ch.stream = nil
	}

	var s = createLms_stream_t()
	s.SetChannel(uint(channelNumber))
	s.SetDataFmt(d.IQFormat)
	s.SetFifoSize(fifoSize)
	s.SetIsTx(!isRX)
	s.SetThroughputVsLatency(0.5)

	if limewrap.LMS_SetupStream(d.dev, s) != 0 {
		panic(fmt.Sprintf("Failed to set stream in %s at %s: %s", d.DeviceInfo.DeviceName, d.DeviceInfo.Media, limewrap.LMS_GetLastErrorMessage()))
	}

	ch.stream = s
}

func (d *LMSDevice) deviceLoop() {

	var cachedActiveChannels = make([]LMSChannel, 0)

	lmsDataChannel := make(chan channelMessage)

	// Check active
	for i := 0; i < len(d.RXChannels); i++ {
		var ch = d.RXChannels[i]
		if ch.stream != nil {
			cachedActiveChannels = append(cachedActiveChannels, *ch)
		}
	}
	// TODO: TX
	//for i := 0; i < len(d.TXChannels); i++ {
	//	var ch = d.TXChannels[i]
	//	if ch.stream != nil {
	//		cachedActiveChannels = append(cachedActiveChannels, ch)
	//		ch.start()
	//	}
	//}

	streamControl := make([]chan bool, len(cachedActiveChannels))

	for i := 0; i < len(cachedActiveChannels); i++ {
		streamControl[i] = make(chan bool)
		ch := cachedActiveChannels[i]
		ch.start()
		go streamLoop(lmsDataChannel, streamControl[i], ch)
	}

	// Notify Main thread that we're done caching
	//log.Println("Device Loop ready.")
	d.controlChan <- true
	//log.Println("Device Loop running with", len(cachedActiveChannels), "channels")
	running := true
	for running {
		select {
		case _ = <-d.controlChan:
			running = false
		case msg := <-lmsDataChannel:
			if d.callback != nil {
				d.callback(msg.data, msg.channel, msg.timestamp)
			}
		}
	}

	// Wait for stopping streams
	//log.Println("Stopping streams")
	for i := 0; i < len(streamControl); i++ {
		select {
		case streamControl[i] <- true: // Send close signal
		case <-lmsDataChannel: // Discard any data received in channel
		}
	}
	d.controlChan <- true
}

// endregion
// region Public Methods
func (d *LMSDevice) SetCallback(cb func([]complex64, int, uint64)) {
	d.callback = cb
}

func (d *LMSDevice) SetGainDB(channelNumber int, isRX bool, gain uint) {
	if limewrap.LMS_SetGaindB(d.dev, !isRX, int64(channelNumber), gain) != 0 {
		panic(fmt.Sprintf("Failed to set channel gain in %s at %s: %s", d.DeviceInfo.DeviceName, d.DeviceInfo.Media, limewrap.LMS_GetLastErrorMessage()))
	}
}

func (d *LMSDevice) SetGainNormalized(channelNumber int, isRX bool, gain float64) {
	if limewrap.LMS_SetNormalizedGain(d.dev, !isRX, int64(channelNumber), gain) != 0 {
		panic(fmt.Sprintf("Failed to set channel gain in %s at %s: %s", d.DeviceInfo.DeviceName, d.DeviceInfo.Media, limewrap.LMS_GetLastErrorMessage()))
	}
}

func (d *LMSDevice) GetGainDB(channelNumber int, isRX bool) (gain uint) {
	if limewrap.LMS_GetGaindB(d.dev, !isRX, int64(channelNumber), &gain) != 0 {
		panic(fmt.Sprintf("Failed to get channel gain in %s at %s: %s", d.DeviceInfo.DeviceName, d.DeviceInfo.Media, limewrap.LMS_GetLastErrorMessage()))
	}
	return gain
}

func (d *LMSDevice) GetGainNormalized(channelNumber int, isRX bool) (gain float64) {
	if limewrap.LMS_GetNormalizedGain(d.dev, !isRX, int64(channelNumber), &gain) != 0 {
		panic(fmt.Sprintf("Failed to get channel gain in %s at %s: %s", d.DeviceInfo.DeviceName, d.DeviceInfo.Media, limewrap.LMS_GetLastErrorMessage()))
	}
	return gain
}

func (d *LMSDevice) GetTemperature() (temp float64) {
	if limewrap.LMS_GetChipTemperature(d.dev, 0, &temp) != 0 {
		panic(fmt.Sprintf("Failed to get chip temperature in %s at %s: %s", d.DeviceInfo.DeviceName, d.DeviceInfo.Media, limewrap.LMS_GetLastErrorMessage()))
	}
	return temp
}

func (d *LMSDevice) SetLPF(channelNumber int, isRX bool, bandwidth float64) {
	if limewrap.LMS_SetLPFBW(d.dev, !isRX, int64(channelNumber), bandwidth) != 0 {
		panic(fmt.Sprintf("Failed to set LPF Bandwidth in %s at %s: %s", d.DeviceInfo.DeviceName, d.DeviceInfo.Media, limewrap.LMS_GetLastErrorMessage()))
	}
}

func (d *LMSDevice) GetLPF(channelNumber int, isRX bool) (bandwidth float64) {
	if limewrap.LMS_GetLPFBW(d.dev, !isRX, int64(channelNumber), &bandwidth) != 0 {
		panic(fmt.Sprintf("Failed to get LPF Bandwidth in %s at %s: %s", d.DeviceInfo.DeviceName, d.DeviceInfo.Media, limewrap.LMS_GetLastErrorMessage()))
	}

	return bandwidth
}

func (d *LMSDevice) EnableLPF(channelNumber int, isRX bool) {
	if limewrap.LMS_SetLPF(d.dev, !isRX, int64(channelNumber), true) != 0 {
		panic(fmt.Sprintf("Failed to enable LPF in %s at %s: %s", d.DeviceInfo.DeviceName, d.DeviceInfo.Media, limewrap.LMS_GetLastErrorMessage()))
	}
}

func (d *LMSDevice) DisableLPF(channelNumber int, isRX bool) {
	if limewrap.LMS_SetLPF(d.dev, !isRX, int64(channelNumber), false) != 0 {
		panic(fmt.Sprintf("Failed to disable LPF in %s at %s: %s", d.DeviceInfo.DeviceName, d.DeviceInfo.Media, limewrap.LMS_GetLastErrorMessage()))
	}
}

func (d *LMSDevice) SetDigitalFilter(channelNumber int, isRX bool, bandwidth float64) {
	var ch *LMSChannel
	if isRX {
		ch = d.RXChannels[channelNumber]
	} else {
		ch = d.TXChannels[channelNumber]
	}

	ch.advancedFiltering = false

	ch.currentDigitalBandwidth = bandwidth

	if ch.currentDigitalBandwidth == 0 {
		panic(fmt.Sprintf("Cannot enable digital filter at channel %d because no bandwidth is set! Call SetDigitalFilter first.", channelNumber))
	}

	limewrap.LMS_SetGFIRLPF(d.dev, isRX, int64(channelNumber), ch.digitalFilterEnabled, ch.currentDigitalBandwidth)
}

func (d *LMSDevice) EnableDigitalFilter(channelNumber int, isRX bool) {
	var ch *LMSChannel
	if isRX {
		ch = d.RXChannels[channelNumber]
	} else {
		ch = d.TXChannels[channelNumber]
	}

	if !ch.advancedFiltering {
		if ch.currentDigitalBandwidth == 0 {
			panic(fmt.Sprintf("Cannot enable digital filter at channel %d because no bandwidth is set! Call SetDigitalFilter first.", channelNumber))
		}

		if limewrap.LMS_SetGFIRLPF(d.dev, isRX, int64(channelNumber), true, ch.currentDigitalBandwidth) != 0 {
			panic(fmt.Sprintf("Failed to enable Digital LPF in %s at %s: %s", d.DeviceInfo.DeviceName, d.DeviceInfo.Media, limewrap.LMS_GetLastErrorMessage()))
		}
	} else {
		panic("Advanced Filtering is enabled. Please use EnableGFIR instead")
	}

	ch.digitalFilterEnabled = true
}

func (d *LMSDevice) DisableDigitalFilter(channelNumber int, isRX bool) {
	var ch *LMSChannel
	if isRX {
		ch = d.RXChannels[channelNumber]
	} else {
		ch = d.TXChannels[channelNumber]
	}

	if limewrap.LMS_SetGFIRLPF(d.dev, isRX, int64(channelNumber), false, ch.currentDigitalBandwidth) != 0 {
		panic(fmt.Sprintf("Failed to disable Digital LPF in %s at %s: %s", d.DeviceInfo.DeviceName, d.DeviceInfo.Media, limewrap.LMS_GetLastErrorMessage()))
	}
	ch.digitalFilterEnabled = false
}

func (d *LMSDevice) EnableChannel(channelNumber int, isRX bool) {
	if limewrap.LMS_EnableChannel(d.dev, !isRX, int64(channelNumber), true) != 0 {
		panic(fmt.Sprintf("Failed to enable channel in %s at %s: %s", d.DeviceInfo.DeviceName, d.DeviceInfo.Media, limewrap.LMS_GetLastErrorMessage()))
	}
	d.setupStream(channelNumber, isRX)
}

func (d *LMSDevice) DisableChannel(channelNumber int, isRX bool) {
	if limewrap.LMS_EnableChannel(d.dev, !isRX, int64(channelNumber), false) != 0 {
		panic(fmt.Sprintf("Failed to disable channel in %s at %s: %s", d.DeviceInfo.DeviceName, d.DeviceInfo.Media, limewrap.LMS_GetLastErrorMessage()))
	}
}

func (d *LMSDevice) SetAntenna(antennaNumber, channelNumber int, isRX bool) {
	if limewrap.LMS_SetAntenna(d.dev, !isRX, int64(channelNumber), int64(antennaNumber)) != 0 {
		panic(fmt.Sprintf("Failed to set antenna in %s at %s: %s", d.DeviceInfo.DeviceName, d.DeviceInfo.Media, limewrap.LMS_GetLastErrorMessage()))
	}
}

func (d *LMSDevice) SetAntennaByName(name string, channelNumber int, isRX bool) {
	var ant *LMSAntenna
	if isRX {
		var c = d.RXChannels[channelNumber]
		for i := 0; i < len(c.Antennas); i++ {
			var a = &c.Antennas[i]
			if strings.ToLower(a.Name) == strings.ToLower(name) {
				ant = a
				break
			}
		}
	} else {
		var c = d.TXChannels[channelNumber]
		for i := 0; i < len(c.Antennas); i++ {
			var a = &c.Antennas[i]
			if strings.ToLower(a.Name) == strings.ToLower(name) {
				ant = a
				break
			}
		}
	}

	if ant == nil {
		panic(fmt.Sprintf("Cannot find antenna with name %s.", name))
	}

	ant.Set()
}

func (d *LMSDevice) Start() {
	if !d.running {
		d.running = true
		go d.deviceLoop()
		//log.Println("Waiting for device loop be ready")
		<-d.controlChan
		//log.Println("Device started")
	} else {
		fmt.Fprintf(os.Stderr, "Device already running")
	}
}

func (d *LMSDevice) Stop() {
	if d.running {
		d.running = false
		d.controlChan <- false
		//log.Println("Waiting loop to stop")
		<-d.controlChan
	} else {
		fmt.Fprintf(os.Stderr, "Device not running")
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

func (d *LMSDevice) SetCenterFrequency(channelNumber int, isRX bool, centerFrequency float64) {
	if limewrap.LMS_SetLOFrequency(d.dev, !isRX, int64(channelNumber), centerFrequency) != 0 {
		panic(fmt.Sprintf("Failed to set Frequency in %s at %s: %s", d.DeviceInfo.DeviceName, d.DeviceInfo.Media, limewrap.LMS_GetLastErrorMessage()))
	}
}

func (d *LMSDevice) GetCenterFrequency(channelNumber int, isRX bool) (centerFrequency float64) {
	if limewrap.LMS_GetLOFrequency(d.dev, !isRX, int64(channelNumber), &centerFrequency) != 0 {
		panic(fmt.Sprintf("Failed to set Frequency in %s at %s: %s", d.DeviceInfo.DeviceName, d.DeviceInfo.Media, limewrap.LMS_GetLastErrorMessage()))
	}
	return centerFrequency
}

func (d *LMSDevice) Close() {
	Close(d)
}

func (d *LMSDevice) String() string {
	var str = fmt.Sprintf("LMSDevice(%s)", d.DeviceInfo.DeviceName)

	str = fmt.Sprintf("%s\nMinimum Sample Rate: %14.0f sps", str, d.MinimumSampleRate)
	str = fmt.Sprintf("%s\nMinimum Sample Rate: %14.0f sps", str, d.MaximumSampleRate)

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

// endregion
