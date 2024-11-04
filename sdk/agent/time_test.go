package agent

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	before := currentTime().Unix()
	time.Sleep(1*time.Second + 1*time.Millisecond)
	now := currentTime().Unix()
	require.True(t, now-before >= 1)
}

// go test -v -run ^$  -bench 'BenchmarkTime' -benchtime=5s -benchmem -cpuprofile profile.pprof
// go tool pprof -http=":8081" profile.pprof
func BenchmarkTime(b *testing.B) {
	b.Run("currentTime", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			currentTime()
		}
	})
	b.Run("time.Now", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			time.Now()
		}
	})
}
