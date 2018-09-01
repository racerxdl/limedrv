package main

import (
	"github.com/racerxdl/limedrv"
	"log"
	"os"
	"time"
)


func OnSamples(data []complex64, channel int, timestamp uint64) {
	log.Println("Received samples from channel", channel, "with timestamp", timestamp)
}

func main() {
	devices := limedrv.GetDevices()

	log.Printf("Found %d devices.\n", len(devices))

	if len(devices) == 0 {
		log.Println("No devices found.")
		os.Exit(1)
	}

	if len(devices) > 1 {
		log.Println("More than one device found. Selecting first one.")
	}

	var di = devices[0]

	log.Printf("Opening device %s\n", di.DeviceName)

	var d = limedrv.Open(di)
	log.Println("Opened!")

	log.Println(d.String())

	d.EnableChannel(limedrv.ChannelA, true)
	d.EnableChannel(limedrv.ChannelB, true)
	d.SetAntennaByName("LNAW", limedrv.ChannelA, true)
	d.SetAntennaByName("LNAW", limedrv.ChannelB, true)

	d.SetCallback(OnSamples)

	d.Start()

	time.Sleep(500 * time.Millisecond)

	d.Stop()

	log.Println("Closing")
	d.Close()

	log.Println("Closed!")
}