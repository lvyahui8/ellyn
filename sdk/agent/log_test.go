package agent

import (
	"github.com/lvyahui8/ellyn/sdk/common/asserts"
	ll "log"
	"os"
	"testing"
)

func TestLogger_Info(t *testing.T) {
	initLogger()
	for i := 0; i < 100000; i++ {
		log.Info("hello world")
	}
	log.file.w.Flush()
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
		f, err := os.OpenFile("logs/test.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		asserts.IsNil(err)
		defer f.Close()
		ll.SetOutput(f)
		for i := 0; i < b.N; i++ {
			ll.Println("This is a test log entry")
		}
	})
}
