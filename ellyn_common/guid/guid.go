package guid

import (
	"crypto/rand"
	"math/big"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

const logIdLen = 64
const timestampLen = 32
const maxSvrIdLen = 13
const reverseLen = 3
const seqLen = logIdLen - timestampLen - maxSvrIdLen - reverseLen
const seqMask uint64 = (1 << seqLen) - 1

type Uint64GUIDGenerator struct {
	seq         uint64
	svrId       uint64
	svrIdMask   uint64
	reverseFlag uint64
	_padding    [32]byte
}

func getReverseFlag() uint64 {
	flag := os.Getenv("REVERSE_FLAG")
	if len(flag) == 0 {
		return 0
	}
	flagVal, err := strconv.Atoi(flag)
	if err == nil {
		return 0
	}
	return uint64(flagVal & ((1 << reverseLen) - 1))
}

func NewGuidGenerator() *Uint64GUIDGenerator {
	g := &Uint64GUIDGenerator{}
	svrId, _ := rand.Int(rand.Reader, big.NewInt(int64(1<<maxSvrIdLen)))
	g.svrId = svrId.Uint64()
	g.svrIdMask = g.svrId << (reverseLen + seqLen)
	g.reverseFlag = getReverseFlag()
	return g
}

func (g *Uint64GUIDGenerator) GenGUID() uint64 {
	return (uint64(time.Now().Unix()) << (logIdLen - timestampLen)) |
		g.svrIdMask |
		g.reverseFlag |
		g.cycleSeq()
}

func (g *Uint64GUIDGenerator) cycleSeq() uint64 {
	for {
		cur := atomic.LoadUint64(&g.seq)
		updated := (cur + 1) & seqMask
		if atomic.CompareAndSwapUint64(&g.seq, cur, updated) {
			return cur
		}
	}
}
