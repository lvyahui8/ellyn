package benchmark

import "testing"

func BenchmarkWrite2DevNull(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Write2DevNull("hello world")
	}
}

func BenchmarkWrite2TmpFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Write2TmpFile("hello world")
	}
}
