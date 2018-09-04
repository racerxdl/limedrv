package limedrv

import "github.com/racerxdl/limedrv/limewrap"

// Preset of channel IDs by name. To be used in channel calls.
const (
	// ChannelA represents the ID of Channel A in LMS Devices ( = 0 )
	ChannelA = 0

	// ChannelB represents the ID of Channel B in LMS Devices ( = 1 )
	ChannelB = 1
)

// Preset of Antenna Names to be used in SetAntennaByName
const (
	// RX Antennas
	LNAW = "LNAW"
	LNAH = "LNAH"
	LNAL = "LNAL"

	// Loopback Antennas (works for both RX and TX)
	LB1 = "LB1"
	LB2 = "LB2"

	// TX Antennas
	BAND1 = "BAND1"
	BAND2 = "BAND2"

	// Not connected
	NONE = "NONE"
)

const fifoSize = 16384 // Samples

// IQ Formats to be set in IQFormat of LMSDevice. This sets the communication between the LMS Device and the computer.
var (
	// FormatFloat32 defines the output of LMS Device to have samples using 32 bit float
	FormatFloat32 = limewrap.Lms_stream_tLMS_FMT_F32
	// FormatInt16 defines the output of LMS Device to have samples using 16 bit int
	FormatInt16 = limewrap.Lms_stream_tLMS_FMT_I16
	// FormatInt12 defines the output of LMS Device to have samples using 12 bit int
	FormatInt12 = limewrap.Lms_stream_tLMS_FMT_I12
)
