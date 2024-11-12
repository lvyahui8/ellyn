package agent

import (
	"github.com/stretchr/testify/require"
	"math"
	"math/rand"
	"testing"
)

func TestNewRandomSampling(t *testing.T) {
	rs := newRandomSampling(0)
	require.Equal(t, uint64(0), rs.target)
	rs = newRandomSampling(1)
	require.Equal(t, uint64(math.MaxUint64), rs.target)
	rs = newRandomSampling(0.5)
	require.True(t, uint64(math.MaxUint64>>1) >= rs.target)
}

func testRate(t *testing.T, targetRate float32) {
	rs := newRandomSampling(targetRate)
	sum := 0
	total := 1000000
	for i := 0; i < total; i++ {
		if rs.hit() {
			sum++
		}
	}
	realRate := float32(sum) / float32(total)
	t.Logf("%f", realRate)
	require.True(t, math.Abs(float64(realRate-targetRate)) <= 0.001) // 误差不超过即可
}

func TestRandomSampling(t *testing.T) {
	testRate(t, 0.5)
	testRate(t, 0.01)
	testRate(t, 0.001)
}

func BenchmarkRandomSampling(b *testing.B) {
	rs := newRandomSampling(1)
	b.Run("samplingRandom", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rs.random()
		}
	})
	rand.Seed(currentTime().UnixMilli())
	b.Run("goRandom", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rand.Uint64()
		}
	})
}
