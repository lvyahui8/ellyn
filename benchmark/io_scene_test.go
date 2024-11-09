package benchmark

import "testing"

const data = `2016-10-25 06:21:25 [Info] ellyn.go site:cn|lang:cn|msg:build success`

func BenchmarkWrite2DevNull(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Write2DevNull(data)
	}
}

// go test -v -run ^$  -bench 'BenchmarkWrite2TmpFile' -benchtime=5s -benchmem -cpuprofile profile.pprof -memprofile memprofile.pprof
// go tool pprof -http=":8081" profile.pprof
// go tool pprof -http=":8082" memprofile.pprof
func BenchmarkWrite2TmpFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Write2TmpFile(data)
	}
}
