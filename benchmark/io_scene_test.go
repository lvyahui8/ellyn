package benchmark

import (
	"testing"
)

const logLine = `2016-10-25 06:21:25 [Info] 20241011230455C8B089895F0389 ellyn.go site:cn|stress:N|lang:cn|msg:build success`

func BenchmarkWrite2DevNull(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Write2DevNull(logLine)
	}
}

// go test -v -run ^$  -bench 'BenchmarkWrite2TmpFile' -benchtime=5s -benchmem -cpuprofile profile.pprof -memprofile memprofile.pprof
// go tool pprof -http=":8081" profile.pprof
// go tool pprof -http=":8082" memprofile.pprof
func BenchmarkWrite2TmpFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Write2TmpFile(logLine)
	}
}

func TestLocalPipeReadWrite(t *testing.T) {
	LocalPipeReadWrite(logLine)
}

func BenchmarkLocalPipeReadWrite(b *testing.B) {
	for i := 0; i < b.N; i++ {
		LocalPipeReadWrite(logLine)
	}
}

func BenchmarkSerialNetRequest(b *testing.B) {
	srv := mockHttpServer()
	defer srv.Close()
	urls := generateLinks(srv.URL)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		syncCrawl(urls)
	}
}

func BenchmarkConcurrentNetRequest(b *testing.B) {
	srv := mockHttpServer()
	defer srv.Close()
	urls := generateLinks(srv.URL)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		concurrentCrawl(urls)
	}
}
