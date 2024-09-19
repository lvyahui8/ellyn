## Ringbuffer性能测试

测试方法：分别使用{1,10,100,500}个生产者/消费者(数量一致)，每个生产者写入100000个元素，测试并发读写性能


```shell
go test -v -run ^$  -bench BenchmarkQueue -benchtime=10s -benchmem
```


> 结论：
> - 当生产者消费者为1时，RingBuffer > channel > linkedBlockQueue
> - 当生产者消费者较小时（<10），RingBuffer >> linkedBlockQueue > channel 
> - 当生产者较大时(>100), RingBuffer >> linkedBlockQueue >> channel
> 
> 并发冲突越多的时候，ringBuffer性能约明显，当生产者消费者=100时，RingBuffer性能是channel的18倍，是linkedBlockQueue的6倍。并且ringBuffer内存分配很少，而linkedBlockQueue内存分配次数很高。

```text
goos: windows
goarch: amd64
pkg: github.com/lvyahui8/ellyn/ellyn_common/collections
cpu: AMD Ryzen 7 4800U with Radeon Graphics
BenchmarkQueue
BenchmarkQueue/1_RingBuffer
BenchmarkQueue/1_RingBuffer-16              1808           6189196 ns/op             192 B/op          7 allocs/op
BenchmarkQueue/1_channelBuffer
BenchmarkQueue/1_channelBuffer-16           1568           8356579 ns/op             189 B/op          7 allocs/op
BenchmarkQueue/1_linkedBlockQueue
BenchmarkQueue/1_linkedBlockQueue-16         990          11683339 ns/op         1600234 B/op     100008 allocs/op
BenchmarkQueue/10_RingBuffer
BenchmarkQueue/10_RingBuffer-16               66         169375852 ns/op            2773 B/op         28 allocs/op
BenchmarkQueue/10_channelBuffer
BenchmarkQueue/10_channelBuffer-16            16         684387012 ns/op            2121 B/op         27 allocs/op
BenchmarkQueue/10_linkedBlockQueue
BenchmarkQueue/10_linkedBlockQueue-16         30         429091237 ns/op        16003789 B/op    1000058 allocs/op
BenchmarkQueue/100_RingBuffer
BenchmarkQueue/100_RingBuffer-16              18         892660172 ns/op           17606 B/op        218 allocs/op
BenchmarkQueue/100_channelBuffer
BenchmarkQueue/100_channelBuffer-16            1        12599217000 ns/op          13664 B/op        230 allocs/op
BenchmarkQueue/100_linkedBlockQueue
BenchmarkQueue/100_linkedBlockQueue-16         3        4411482867 ns/op        160070496 B/op  10000905 allocs/op
BenchmarkQueue/500_RingBuffer
BenchmarkQueue/500_RingBuffer-16               2        8278503150 ns/op          100416 B/op       1087 allocs/op
BenchmarkQueue/500_channelBuffer
BenchmarkQueue/500_channelBuffer-16            1        83803768800 ns/op          65568 B/op       1104 allocs/op
BenchmarkQueue/500_linkedBlockQueue
BenchmarkQueue/500_linkedBlockQueue-16         1        21573212700 ns/op       800390176 B/op  50004902 allocs/op
PASS
ok      github.com/lvyahui8/ellyn/ellyn_common/collections      258.178s

```