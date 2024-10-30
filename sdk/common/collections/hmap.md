## map 性能测试结果

测试方法：随机生成1000w个完全不相同的数字k，k<1亿 ，打乱顺序让多协程并发写入、读取、删除，每个协程处理其中的一部分号段

> 结论: concurrentMap 在分段为2048时，性能和内存分配明显好于sync.Map。 但在资源竞争率高的时候，性能差异不明显，concurrentMap略好于syncMap。

```shell
go test -v -run ^$  -bench BenchmarkMap -benchtime=5s -benchmem
```

```text
goos: windows
goarch: amd64
pkg: github.com/lvyahui8/ellyn/sdk/common/collections
cpu: AMD Ryzen 7 4800U with Radeon Graphics
BenchmarkMap
BenchmarkMap/concurrentMap_put
BenchmarkMap/concurrentMap_put-16                     18         342347494 ns/op        80002544 B/op    9999784 allocs/op
BenchmarkMap/concurrentMap_read
BenchmarkMap/concurrentMap_read-16                    18         339383994 ns/op        80001507 B/op    9999781 allocs/op
BenchmarkMap/concurrentMap_delete
BenchmarkMap/concurrentMap_delete-16                  61          99936639 ns/op        80000277 B/op    9999779 allocs/op
BenchmarkMap/concurrentMap_R&W
BenchmarkMap/concurrentMap_R&W-16                      9         723391478 ns/op        160000385 B/op  19999569 allocs/op
BenchmarkMap/syncMap_put
BenchmarkMap/syncMap_put-16                            1        8409727400 ns/op        1470761288 B/op 40306906 allocs/op
BenchmarkMap/syncMap_read
BenchmarkMap/syncMap_read-16                           1        6106719400 ns/op        80012880 B/op    9999918 allocs/op
BenchmarkMap/syncMap_delete
BenchmarkMap/syncMap_delete-16                        15         369030780 ns/op        79999968 B/op    9999778 allocs/op
BenchmarkMap/syncMap_R&W
BenchmarkMap/syncMap_R&W-16                            7        1029901486 ns/op        319999952 B/op  29999555 allocs/op
PASS
ok      github.com/lvyahui8/ellyn/sdk/common/collections      70.109s
```
