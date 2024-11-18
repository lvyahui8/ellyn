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

func TestDate(t *testing.T) {
	t.Log(date())
	require.Equal(t, 20241101, getDate(time.UnixMilli(1730475387000)))
	require.Equal(t, 20241230, getDate(time.UnixMilli(1735572987000)))
	require.Equal(t, 20240229, getDate(time.UnixMilli(1709220987000)))
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
