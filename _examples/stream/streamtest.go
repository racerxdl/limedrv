package main

import (
	"github.com/racerxdl/limedrv"
	"log"
	"os"
)


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

	log.Println("Closing")
	d.Close()

	log.Println("Closed!")
}