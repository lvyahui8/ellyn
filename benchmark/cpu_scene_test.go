package benchmark

import "testing"

func TestQuickSort(t *testing.T) {
	for i := 0; i < 100000; i++ {
		arr := []int{4, 5, 1, 7, 8, 10}
		quickSort(arr, 0, len(arr)-1)
	}
}

// go test -v -run ^$  -bench 'BenchmarkQuickSort' -benchtime=5s -benchmem -cpuprofile profile.pprof -memprofile memprofile.pprof
// go tool pprof -http=":8081" profile.pprof
// go tool pprof -http=":8082" memprofile.pprof
func BenchmarkQuickSort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		arr := []int{4, 5, 1, 7, 8, 10}
		quickSort(arr, 0, len(arr)-1)
	}
}
