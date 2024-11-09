package agent

import (
	"testing"
	"unsafe"
)

func TestGraphSize(t *testing.T) {
	t.Log(unsafe.Sizeof(graph{}))
	t.Log(unsafe.Sizeof(graph{}.nodes))
	t.Log(unsafe.Sizeof(graph{}.edges))
}

// go test -v -run ^$  -bench 'BenchmarkGraphAddEdge' -benchtime=5s -benchmem -cpuprofile profile.pprof -memprofile memprofile.pprof
// go tool pprof -http=":8081" profile.pprof
// go tool pprof -http=":8082" memprofile.pprof
func BenchmarkGraphAddEdge(b *testing.B) {
	g := newGraph(0)
	for i := 0; i < b.N; i++ {
		g.addEdge(2, 0)
		g.addEdge(2, 2)
		for k := range g.edges {
			delete(g.edges, k)
		}
	}
}
