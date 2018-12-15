// limedrv is a LMS7 API Wrapper to Go made to be easy to use.
// Currently this documentation is WIP. Some examples are available
// in _examples folder.
package limedrv

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/racerxdl/limedrv/limewrap"
	"runtime"
	"unsafe"
)

type i_deviceinfo struct {
	DeviceName          [64]byte
	FirmwareVersion     [16]byte
	HardwareVersion     [16]byte
	ProtocolVersion     [16]byte
	BoardSerialNumber   uint64
	GatewareVersion     [16]byte
	GatewareTargetBoard [32]byte
}

// DeviceInfo is a struct with driver information required to open a connection
type DeviceInfo struct {
	DeviceName          string
	Media               string
	Module              string
	Addr                string
	Serial              string
	ProtocolVersion     string
	FirmwareVersion     string
	HardwareVersion     string
	GatewareVersion     string
	GatewareTargetBoard string
	origDevInfo         i_deviceinfo
}

func (d *i_deviceinfo) toOrigDevString() string {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, d)
	if err != nil {
		panic(err)
	}

	return buf.String()
}

// GetDevices return an array of available devices in the LMS7 driver.
func GetDevices() []DeviceInfo {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	devCount := limewrap.LMS_GetDeviceList(nil)
	ret := make([]DeviceInfo, devCount)

	if devCount > 0 {
		var z [128]i_deviceinfo
		t := (*string)(unsafe.Pointer(&z))
		limewrap.LMS_GetDeviceList(t)
		for i := 0; i < devCount; i++ {
			ret[i] = idev2dev(z[i])
		}
	}

	return ret
}

// Open opens a device specified by a DeviceInfo instance and returns a reference to LMSDevice
func Open(device DeviceInfo) *LMSDevice {
	var ret = LMSDevice{
		DeviceInfo:  device,
		IQFormat:    FormatInt16,
		controlChan: make(chan bool),
	}

	ret.Advanced = LMSDeviceAdvanced{}

	var origString = device.origDevInfo.toOrigDevString()

	ptr := uintptr(0)

	runtime.LockOSThread()
	v := limewrap.LMS_Open(&ptr, origString, 0)
	runtime.UnlockOSThread()
	ret.dev = ptr

	if v != 0 {
		panic(fmt.Sprintf("Failed to open %s at %s.", device.DeviceName, device.Media))
	}

	ret.init()

	return &ret
}

// Close closes a LMSDevice. This makes the LMSDevice instance useless.
func Close(device *LMSDevice) {
	if limewrap.LMS_Close(device.dev) != 0 {
		panic(fmt.Sprintf("Failed to close %s at %s.", device.DeviceInfo.DeviceName, device.DeviceInfo.Media))
	} else {
		device.dev = 0
	}
}
