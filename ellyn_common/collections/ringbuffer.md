## Ringbuffer性能测试

```shell
go test -v -run ^$  -bench BenchmarkRingBuffer -benchtime=5s -benchmem
```

```text
goos: windows
goarch: amd64
pkg: github.com/lvyahui8/ellyn/ellyn_common/collections
cpu: AMD Ryzen 7 4800U with Radeon Graphics
BenchmarkRingBuffer10000
BenchmarkRingBuffer10000-16          382          16506136 ns/op           26062 B/op        376 allocs/op
BenchmarkRingBuffer
BenchmarkRingBuffer/ringBuffer_readWrite_1
BenchmarkRingBuffer/ringBuffer_readWrite_1-16             577216              9874 ns/op            1296 B/op         33 allocs/op
BenchmarkRingBuffer/channelBuffer_readWrite_1
BenchmarkRingBuffer/channelBuffer_readWrite_1-16          517130             11760 ns/op            1296 B/op         33 allocs/op
BenchmarkRingBuffer/ringBuffer_readWrite_10
BenchmarkRingBuffer/ringBuffer_readWrite_10-16            224670             26767 ns/op            1296 B/op         33 allocs/op
BenchmarkRingBuffer/channelBuffer_readWrite_10
BenchmarkRingBuffer/channelBuffer_readWrite_10-16         168636             35749 ns/op            1296 B/op         33 allocs/op
BenchmarkRingBuffer/ringBuffer_readWrite_100
BenchmarkRingBuffer/ringBuffer_readWrite_100-16            31714            189755 ns/op            1296 B/op         33 allocs/op
BenchmarkRingBuffer/channelBuffer_readWrite_100
BenchmarkRingBuffer/channelBuffer_readWrite_100-16          7354            781929 ns/op            1296 B/op         33 allocs/op
BenchmarkRingBuffer/ringBuffer_readWrite_1000
BenchmarkRingBuffer/ringBuffer_readWrite_1000-16            3050           1942901 ns/op            1296 B/op         33 allocs/op
BenchmarkRingBuffer/channelBuffer_readWrite_1000
BenchmarkRingBuffer/channelBuffer_readWrite_1000-16          625           9612621 ns/op            1296 B/op         33 allocs/op
BenchmarkRingBuffer/ringBuffer_readWrite_10000
BenchmarkRingBuffer/ringBuffer_readWrite_10000-16            332          17939649 ns/op            1296 B/op         33 allocs/op
BenchmarkRingBuffer/channelBuffer_readWrite_10000
BenchmarkRingBuffer/channelBuffer_readWrite_10000-16          55         104589107 ns/op            1296 B/op         33 allocs/op
PASS
ok      github.com/lvyahui8/ellyn/ellyn_common/collections      73.444s

```