package limedrv

import "github.com/racerxdl/limedrv/limewrap"

const ChannelA = 0
const ChannelB = 1

const LNAW = "LNAW"
const LNAH = "LNAH"
const LNAL = "LNAL"
const NONE = "NONE"

const LB1 = "LB1"
const LB2 = "LB2"

const BAND1 = "BAND1"
const BAND2 = "BAND2"

const fifoSize = 16384 // Samples

var FormatFloat32 = limewrap.Lms_stream_tLMS_FMT_F32
var FormatInt16 = limewrap.Lms_stream_tLMS_FMT_I16
var FormatInt12 = limewrap.Lms_stream_tLMS_FMT_I12
