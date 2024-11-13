package benchmark

import (
	"crypto/rand"
	"github.com/lvyahui8/ellyn/api"
	"testing"
	"time"
)

import _ "github.com/lvyahui8/ellyn/api"
import _ "github.com/lvyahui8/ellyn"

func TestQuickSort(t *testing.T) {
	for i := 0; i < 1000000; i++ {
		arr := []int{4, 5, 1, 7, 8, 10}
		quickSort(arr, 0, len(arr)-1)
	}
	time.Sleep(100 * time.Millisecond)
	t.Log(api.Agent.GetGraphCnt()) // 当前采样率0.001，则累计1k
}

func newUnsortedArr() (arr []int) {
	for i := 0; i < 200; i++ {
		arr = append(arr, i)
	}
	return shuffle(arr)
}

func newSortedArr() (arr []int) {
	for i := 0; i < 200; i++ {
		arr = append(arr, i)
	}
	return
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

func BenchmarkBinarySearch(b *testing.B) {
	arr := newSortedArr()
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
	arr := newUnsortedArr()
	for i := 0; i < b.N; i++ {
		shuffle(arr)
	}
}

func BenchmarkStringCompress(b *testing.B) {
	// 4KB 字符串压缩
	data := make([]byte, 4*1024)
	_, _ = rand.Read(data)
	content := string(data)
	for i := 0; i < b.N; i++ {
		StringCompress(content)
	}
}

func BenchmarkEncryptAndDecrypt(b *testing.B) {
	// 100个字符加解密

	for i := 0; i < b.N; i++ {
		EncryptAndDecrypt(logLine)
	}
}
