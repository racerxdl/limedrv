[![Build Status](https://api.travis-ci.org/racerxdl/limedrv.svg?branch=master)](https://travis-ci.org/racerxdl/limedrv) [![Apache License](https://img.shields.io/badge/license-Apache-blue.svg)](https://tldrlegal.com/license/apache-license-2.0-(apache-2.0)) [![Go Report](https://goreportcard.com/badge/github.com/racerxdl/limedrv)](https://goreportcard.com/report/github.com/racerxdl/limedrv)

# limedrv
LimeSuite Wrapper on Go (Driver for LimeSDR Devices)

# Usage

So far I need to do all the comments for the methods (since go auto-generates the documentation).
But while I do that, you can check the examples. The documentation is available at: [https://godoc.org/github.com/racerxdl/limedrv](https://godoc.org/github.com/racerxdl/limedrv)


# Examples

So far there is a functional WBFM Radio that uses SegDSP for demodulating. You can check it at `_examples/limefm`. To compile, just go to the folder and run:

```bash
go build
```

It will generate a `limefm` executable in the folder. It outputs the raw Float32 audio into stdout. For example, you can listen to the radio by using ffplay:

```bash
./limefm -antenna LNAL -centerFrequency 106300000 -channel 0 -gain 0.5 -outputRate 48000 | ffplay -f f32le -ar 48k -ac 1 -
```
