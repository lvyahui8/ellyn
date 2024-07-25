## stack 性能测试结果

> 结论：Push相同元素时性能很好，其他场景与非压缩的差不多。

```shell
go test -v -run ^$  -bench BenchmarkStack -benchtime=5s -benchmem
```

```text
goos: windows
goarch: amd64
pkg: github.com/lvyahui8/ellyn/ellyn_common/collections
cpu: AMD Ryzen 7 4800U with Radeon Graphics
BenchmarkStack
BenchmarkStack/UnsafeStack_Push
BenchmarkStack/UnsafeStack_Push-16              42489004               149.8 ns/op            55 B/op          1 allocs/op
BenchmarkStack/UnsafeStack_Top
BenchmarkStack/UnsafeStack_Top-16               1000000000               2.398 ns/op           0 B/op          0 allocs/op
BenchmarkStack/UnsafeStack_Pop
BenchmarkStack/UnsafeStack_Pop-16               1000000000               2.583 ns/op           0 B/op          0 allocs/op
BenchmarkStack/UnsafeStack_Push_same
BenchmarkStack/UnsafeStack_Push_same-16         76058860                99.86 ns/op           48 B/op          1 allocs/op
BenchmarkStack/UnsafeStack_Top_same
BenchmarkStack/UnsafeStack_Top_same-16          1000000000               2.406 ns/op           0 B/op          0 allocs/op
BenchmarkStack/UnsafeStack_Pop_same
BenchmarkStack/UnsafeStack_Pop_same-16          1000000000               2.295 ns/op           0 B/op          0 allocs/op
BenchmarkStack/UnsafeCompressedStack_Push
BenchmarkStack/UnsafeCompressedStack_Push-16            44350060               125.6 ns/op            87 B/op          2 allocs/op
BenchmarkStack/UnsafeCompressedStack_Top
BenchmarkStack/UnsafeCompressedStack_Top-16             1000000000               2.505 ns/op           0 B/op          0 allocs/op
BenchmarkStack/UnsafeCompressedStack_Pop
BenchmarkStack/UnsafeCompressedStack_Pop-16             1000000000               2.111 ns/op           0 B/op          0 allocs/op
BenchmarkStack/UnsafeCompressedStack_Push_same
BenchmarkStack/UnsafeCompressedStack_Push_same-16       1000000000               5.569 ns/op           0 B/op          0 allocs/op
BenchmarkStack/UnsafeCompressedStack_Top_same
BenchmarkStack/UnsafeCompressedStack_Top_same-16        1000000000               2.566 ns/op           0 B/op          0 allocs/op
BenchmarkStack/UnsafeCompressedStack_Pop_same
BenchmarkStack/UnsafeCompressedStack_Pop_same-16        1000000000               3.254 ns/op           0 B/op          0 allocs/op
PASS
ok      github.com/lvyahui8/ellyn/ellyn_common/collections      174.620s
```