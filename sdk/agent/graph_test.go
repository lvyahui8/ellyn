package agent

import "testing"

// go test -v -run ^$  -bench 'BenchmarkGraphAddEdge' -benchtime=5s -benchmem -cpuprofile profile.pprof -memprofile memprofile.pprof
// go tool pprof -http=":8081" profile.pprof
// go tool pprof -http=":8082" memprofile.pprof
func BenchmarkGraphAddEdge(b *testing.B) {
	g := newGraph(0)
	for i := 0; i < b.N; i++ {
		g.addEdge(4, 3)
	}
}
