package agent

import (
	"github.com/lvyahui8/ellyn/sdk/common/ctime"
	"math"
)

var sampling *randomSampling

func initSampling() {
	sampling = newRandomSampling(float32(conf.SamplingRate))
}

type randomSampling struct {
	target uint64
	cur    uint64
}

func newRandomSampling(samplingRate float32) *randomSampling {
	if samplingRate <= 0 {
		samplingRate = 0
	}
	if samplingRate >= 1 {
		samplingRate = 1
	}
	rs := &randomSampling{}
	precision := uint64(100000000)
	if samplingRate == 1 {
		rs.target = math.MaxUint64
	} else {
		rs.target = math.MaxUint64 / precision * (uint64(float32(precision) * samplingRate))
	}
	rs.cur = uint64(ctime.CurrentTime().UnixMicro())
	return rs
}

func (rs *randomSampling) hit() bool {
	return rs.random() < rs.target
}

func (rs *randomSampling) random() uint64 {
	x := rs.cur
	if x == 0 {
		x = uint64(ctime.CurrentTime().UnixMicro())
	}
	x ^= x << 13
	x ^= x >> 7
	x ^= x << 17
	rs.cur = x
	return x
}
