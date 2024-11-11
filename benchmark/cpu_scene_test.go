package benchmark

import (
	"github.com/lvyahui8/ellyn/api"
	"testing"
	"time"
)

import _ "github.com/lvyahui8/ellyn/api"
import _ "github.com/lvyahui8/ellyn"

func TestQuickSort(t *testing.T) {
	for i := 0; i < 100000; i++ {
		arr := []int{4, 5, 1, 7, 8, 10}
		quickSort(arr, 0, len(arr)-1)
	}
	time.Sleep(100 * time.Millisecond)
	t.Log(api.Agent.GetGraphCnt()) // 当前采样率0.01，则累计1w
}

func newUnsortedArr() []int {
	return []int{4, 5, 1, 7, 8, 10, 11, -1, -9, 199, 29, 10, 14, 28}
}

// go test -v -run ^$  -bench 'BenchmarkQuickSort' -benchtime=5s -benchmem -cpuprofile profile.pprof -memprofile memprofile.pprof
// go tool pprof -http=":8081" profile.pprof
// go tool pprof -http=":8082" memprofile.pprof
func BenchmarkQuickSort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		arr := newUnsortedArr()
		quickSort(arr, 0, len(arr)-1)
	}
}

func BenchmarkBinarySerch(b *testing.B) {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	for i := 0; i < b.N; i++ {
		_ = binarySearch(arr, 14)
	}
}

func BenchmarkBubbleSort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		arr := newUnsortedArr()
		bubbleSort(arr)
	}
}

func BenchmarkShuffle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		arr := newUnsortedArr()
		shuffle(arr)
	}
}
