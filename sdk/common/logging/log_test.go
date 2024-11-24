package logging

import (
	"fmt"
	"github.com/lvyahui8/ellyn/sdk/common/asserts"
	"github.com/lvyahui8/ellyn/sdk/common/ctime"
	"github.com/stretchr/testify/require"
	ll "log"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestLogFileDumpName(t *testing.T) {
	initLogger()
	file := log.file.getBaseLogFile()
	name := log.file.dumpFileName()
	require.True(t, strings.HasPrefix(name, file))
	//require.True(t, strings.HasSuffix(name, ".0"))
	idx, err := strconv.Atoi(name[len(file)+len(fmt.Sprintf(".%d.", ctime.Date())):])
	require.Nil(t, err)
	require.True(t, idx >= 0)
}

func TestLogger_Info(t *testing.T) {
	initLogger()
	for i := 0; i < 200000; i++ {
		log.Info("hello world")
	}
	time.Sleep(1 * time.Second)
	log.flush()
}

func TestLogger_InfoKV(t *testing.T) {
	initLogger()
	log.InfoKV(Empty().Str("name", "yah").Int("age", 1))
	time.Sleep(1 * time.Second)
	log.flush()
}

// go test -v -run ^$  -bench 'BenchmarkLogger/asyncLogger' -benchtime=5s -benchmem -cpuprofile profile.pprof -memprofile memprofile.pprof
// go tool pprof -http=":8081" profile.pprof
// go tool pprof -http=":8082" memprofile.pprof
func BenchmarkLogger(b *testing.B) {
	initLogger()
	b.Run("asyncLogger", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			log.Info("hello world")
		}
	})
	b.Run("syncLogger", func(b *testing.B) {
		f, err := os.OpenFile("logs/test.logging", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		asserts.IsNil(err)
		defer f.Close()
		ll.SetOutput(f)
		for i := 0; i < b.N; i++ {
			ll.Println("This is a test logging entry")
		}
	})
}

// go test -v -run ^$  -bench 'BenchmarkLoggerFormatAndKV/kvLog' -benchtime=5s -benchmem -cpuprofile profile.pprof -memprofile memprofile.pprof
// go tool pprof -http=":8081" profile.pprof
// go tool pprof -http=":8082" memprofile.pprof
func BenchmarkLoggerFormatAndKV(b *testing.B) {
	initLogger()
	b.Run("formatLog", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			log.Info("name:%s|age:%d|suc:%b", "yah", 1, true)
		}
	})
	b.Run("kvLog", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			log.InfoKV(Empty().Str("name", "yah").Int("age", 1).Bool("suc", true))
		}
	})
}
