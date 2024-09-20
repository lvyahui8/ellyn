### RingBuffer VS Map写性能对比

目的：验证Pop方法构建调用链路用RingBuffer异步实现还是用map实现，哪一个性能影响最小

测试方法：对比以下三种场景

- ringBuffer: 往RingBuffer里写值1（队列里写任何值没区别）
- mapSeqWrite: 往map里面写递增值, 值一直递增直到math.MaxInt
- mapEachWrite: 往map里面写循环写值，值从0到255循环写
- mapNormalWrite: 往已经有1000个元素（0<=key<1000）的固定写500

> 结论： **RingBuffer性能要好于map**，在map元素较少时，RingBuffer的性能是map的4-5倍。而map随着元素的增加写性能会明显降低，最后可能跟RingBuffer相差几个数量级。
> 但这里还需要考虑2点
> - 这里只验证了RingBuffer单线程写的性能，RingBuffer并发写的性能可能会有所降低
> - 多数时候，一次调用，走到的方法不会太多（仅限有采集的方法），一般不会超过100个。这时候RingBuffer并没有非常明显的优势

```shell
go test -v -run ^$  -bench BenchmarkRingBufferAndMap -benchtime=10s -benchmem
```

```text
goos: windows
goarch: amd64
pkg: github.com/lvyahui8/ellyn/ellyn_common/collections
cpu: AMD Ryzen 7 4800U with Radeon Graphics
BenchmarkRingBufferAndMap
BenchmarkRingBufferAndMap/ringBuffer
BenchmarkRingBufferAndMap/ringBuffer-16                 1000000000               2.188 ns/op           0 B/op          0 allocs/op
BenchmarkRingBufferAndMap/mapSeqWrite
BenchmarkRingBufferAndMap/mapSeqWrite-16                77492746               237.8 ns/op            41 B/op          0 allocs/op
BenchmarkRingBufferAndMap/mapEachWrite
BenchmarkRingBufferAndMap/mapEachWrite-16               1000000000               9.690 ns/op           0 B/op          0 allocs/op
BenchmarkRingBufferAndMap/mapNormalWrite
BenchmarkRingBufferAndMap/mapNormalWrite-16             1000000000               8.359 ns/op           0 B/op          0 allocs/op
PASS
ok      github.com/lvyahui8/ellyn/ellyn_common/collections      41.103s
```