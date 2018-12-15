package limedrv

import "C"
import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/racerxdl/limedrv/limewrap"
	"runtime"
	"strings"
	"unsafe"
)

func cleanString(s string) string {
	return strings.Trim(s, "\u0000 ")
}

type channelMessage struct {
	channel   int
	data      []complex64
	timestamp uint64
}

func streamLoop(c chan<- channelMessage, con chan bool, channel LMSChannel) {
	var err error
	//fmt.Fprintf(os.Stderr,"Worker Started")
	running := true
	sampleLength := 4
	if channel.parent.IQFormat == FormatInt16 || channel.parent.IQFormat == FormatInt12 {
		sampleLength = 2
	}
	buff := make([]byte, fifoSize*sampleLength*2) // 16k IQ samples
	zeroPointer := uintptr(unsafe.Pointer(&buff[0]))

	m := limewrap.NewLms_stream_meta_t()
	m.SetTimestamp(0)
	m.SetFlushPartialPacket(false)
	m.SetWaitForTimestamp(false)
	//fmt.Fprintf(os.Stderr,"Worker Running")
	for running {
		select {
		case _ = <-con:
			//fmt.Fprintf(os.Stderr,"Worker Received stop", b)
			running = false
			return
		default:
		}

		recvSamples := limewrap.LMS_RecvStream(channel.stream, zeroPointer, 16384, m, 100)
		if recvSamples > 0 {
			chunk := buff[:sampleLength*recvSamples*2]
			rbuf := bytes.NewReader(chunk)
			cm := channelMessage{
				channel:   channel.parentIndex,
				data:      make([]complex64, recvSamples),
				timestamp: m.GetTimestamp(),
			}

			if sampleLength == 4 {
				// Float32
				v := make([]float32, recvSamples)
				err = binary.Read(rbuf, binary.LittleEndian, &v)
				if err != nil {
					panic(err)
				}
				for i := 0; i < recvSamples; i++ {
					cm.data[i] = complex(v[i*2], v[i*2+1])
				}
			} else {
				// Int16
				//var i16a, i16b int16
				var i16buff = make([]int16, recvSamples*2)
				err = binary.Read(rbuf, binary.LittleEndian, &i16buff)
				for i := 0; i < recvSamples; i++ {
					cm.data[i] = complex(float32(i16buff[i*2])/32768, float32(i16buff[i*2+1])/32768)
				}
			}

			c <- cm
		} else if recvSamples == -1 {
			fmt.Printf("Error receiving samples from channel %d\n", channel.parentIndex)
		}
		runtime.Gosched()
	}
}

func createLms_range_t() limewrap.Lms_range_t {
	return limewrap.NewLms_range_t()
}

func createLms_stream_t() limewrap.Lms_stream_t {
	return limewrap.NewLms_stream_t()
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
			case "media":
				Media = cleanString(k[1])
				break
			case "module":
				Module = cleanString(k[1])
				break
			case "addr":
				Addr = cleanString(k[1])
				break
			case "serial":
				Serial = cleanString(k[1])
				break
			}
		}
	}

	return DeviceInfo{
		DeviceName:          DeviceName,
		Media:               Media,
		Module:              Module,
		Addr:                Addr,
		Serial:              Serial,
		FirmwareVersion:     cleanString(string(deviceinfo.FirmwareVersion[:16])),
		HardwareVersion:     cleanString(string(deviceinfo.HardwareVersion[:16])),
		GatewareVersion:     cleanString(string(deviceinfo.GatewareVersion[:16])),
		GatewareTargetBoard: cleanString(string(deviceinfo.GatewareTargetBoard[:16])),
		origDevInfo:         deviceinfo,
	}
}
