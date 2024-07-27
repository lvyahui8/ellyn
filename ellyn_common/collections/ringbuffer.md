## Ringbuffer性能测试

```shell
go test -v -run ^$  -bench BenchmarkRingBuffer -benchtime=5s -benchmem
```

```text
goos: windows
goarch: amd64
pkg: github.com/lvyahui8/ellyn/ellyn_common/collections
cpu: AMD Ryzen 7 4800U with Radeon Graphics
BenchmarkRingBuffer
BenchmarkRingBuffer/ringBuffer_readWrite_1
BenchmarkRingBuffer/ringBuffer_readWrite_1-16             539079             10520 ns/op            1313 B/op         33 allocs/op
BenchmarkRingBuffer/channelBuffer_readWrite_1
BenchmarkRingBuffer/channelBuffer_readWrite_1-16          489321             12805 ns/op            1299 B/op         33 allocs/op
BenchmarkRingBuffer/ringBuffer_readWrite_10
BenchmarkRingBuffer/ringBuffer_readWrite_10-16            211220             27978 ns/op            1340 B/op         33 allocs/op
BenchmarkRingBuffer/channelBuffer_readWrite_10
BenchmarkRingBuffer/channelBuffer_readWrite_10-16         152728             39356 ns/op            1306 B/op         33 allocs/op
BenchmarkRingBuffer/ringBuffer_readWrite_100
BenchmarkRingBuffer/ringBuffer_readWrite_100-16            29646            204086 ns/op            1614 B/op         37 allocs/op
BenchmarkRingBuffer/channelBuffer_readWrite_100
BenchmarkRingBuffer/channelBuffer_readWrite_100-16          7141            800233 ns/op            1520 B/op         33 allocs/op
BenchmarkRingBuffer/ringBuffer_readWrite_1000
BenchmarkRingBuffer/ringBuffer_readWrite_1000-16            2828           2113992 ns/op            4633 B/op         79 allocs/op
BenchmarkRingBuffer/channelBuffer_readWrite_1000
BenchmarkRingBuffer/channelBuffer_readWrite_1000-16          608           8629558 ns/op            3937 B/op         33 allocs/op
BenchmarkRingBuffer/ringBuffer_readWrite_10000
BenchmarkRingBuffer/ringBuffer_readWrite_10000-16            308          18802652 ns/op           31936 B/op        458 allocs/op
BenchmarkRingBuffer/channelBuffer_readWrite_10000
BenchmarkRingBuffer/channelBuffer_readWrite_10000-16          79         104347427 ns/op           21622 B/op         33 allocs/op
PASS
ok      github.com/lvyahui8/ellyn/ellyn_common/collections      67.490s

```