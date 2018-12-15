package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/racerxdl/limedrv"
	"github.com/racerxdl/segdsp/demodcore"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const sampleRate = 2e6

var demod *demodcore.FMDemod

var output *os.File

var outputRate = flag.Int("outputRate", 48000, "Output Rate in Hz")
var centerFrequency = flag.Float64("centerFrequency", 93.7e6, "Preset for Demodulator Params")
var gain = flag.Float64("gain", 0.5, "Normalized Gain [0-1]")
var antenna = flag.String("antenna", "LNAL", "Antenna Name [LNAL, LNAH, LNAW]")
var channel = flag.Int("channel", 0, "Channel Number [0 => A, 1 => B]")

func OnSamples(data []complex64, _ int, _ uint64) {
	var demodDataI = demod.Work(data)

	if demodDataI != nil {
		demodData := demodDataI.(demodcore.DemodData)
		if output != nil {
			binary.Write(output, binary.LittleEndian, demodData.Data)
		}
	}
}

func main() {

	flag.Parse()

	demod = demodcore.MakeWBFMDemodulator(sampleRate, 120e3, uint32(*outputRate))

	devices := limedrv.GetDevices()

	fmt.Fprintf(os.Stderr,"Found %d devices.\n", len(devices))

	if len(devices) == 0 {
		fmt.Fprintf(os.Stderr,"No devices found.\n")
		os.Exit(1)
	}

	if len(devices) > 1 {
		fmt.Fprintf(os.Stderr,"More than one device found. Selecting first one.\n")
	}

	var di = devices[0]

	fmt.Fprintf(os.Stderr,"Opening device %s\n", di.DeviceName)

	var d = limedrv.Open(di)

	d.SetSampleRate(sampleRate, 8)

	//log.Println(d.String())

	output = os.Stdout

	var ch = d.RXChannels[*channel]

	ch.Enable().
		SetAntennaByName(*antenna).
		SetGainNormalized(*gain).
		SetLPF(1.5e6).
		EnableLPF().
		SetDigitalLPF(300e3).
		EnableDigitalLPF().
		SetCenterFrequency(*centerFrequency)

	d.SetCallback(OnSamples)

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println(sig)
		done <- true
	}()


	d.Start()

	<-done

	d.Stop()

	log.Println("Closing")
	d.Close()

	log.Println("Closed!")
}