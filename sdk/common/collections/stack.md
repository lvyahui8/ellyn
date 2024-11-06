## stack 性能测试结果

> 结论：UnsafeUint32Stack读写性能最优

```shell
go test -v -run ^$  -bench BenchmarkStack -benchtime=5s -benchmem
```

```text
goos: windows
goarch: amd64
pkg: github.com/lvyahui8/ellyn/sdk/common/collections
cpu: AMD Ryzen 7 4800U with Radeon Graphics
BenchmarkStack
BenchmarkStack/UnsafeStack_Push
BenchmarkStack/UnsafeStack_Push-16              34034480               229.8 ns/op            52 B/op          1 allocs/op
BenchmarkStack/UnsafeStack_Top
BenchmarkStack/UnsafeStack_Top-16               1000000000               3.275 ns/op           0 B/op          0 allocs/op
BenchmarkStack/UnsafeStack_Pop
BenchmarkStack/UnsafeStack_Pop-16               1000000000               3.145 ns/op           0 B/op          0 allocs/op
BenchmarkStack/UnsafeStack_Push_same
BenchmarkStack/UnsafeStack_Push_same-16         43446966               142.3 ns/op            48 B/op          1 allocs/op
BenchmarkStack/UnsafeStack_Top_same
BenchmarkStack/UnsafeStack_Top_same-16          1000000000               2.887 ns/op           0 B/op          0 allocs/op
BenchmarkStack/UnsafeStack_Pop_same
BenchmarkStack/UnsafeStack_Pop_same-16          1000000000               4.067 ns/op           0 B/op          0 allocs/op
BenchmarkStack/UnsafeCompressedStack_Push
BenchmarkStack/UnsafeCompressedStack_Push-16            32847302               291.6 ns/op            67 B/op          2 allocs/op
BenchmarkStack/UnsafeCompressedStack_Top
BenchmarkStack/UnsafeCompressedStack_Top-16             1000000000               3.234 ns/op           0 B/op          0 allocs/op
BenchmarkStack/UnsafeCompressedStack_Pop
BenchmarkStack/UnsafeCompressedStack_Pop-16             1000000000               3.575 ns/op           0 B/op          0 allocs/op
BenchmarkStack/UnsafeCompressedStack_Push_same
BenchmarkStack/UnsafeCompressedStack_Push_same-16       502593381               10.71 ns/op            0 B/op          0 allocs/op
BenchmarkStack/UnsafeCompressedStack_Top_same
BenchmarkStack/UnsafeCompressedStack_Top_same-16        1000000000               3.206 ns/op           0 B/op          0 allocs/op
BenchmarkStack/UnsafeCompressedStack_Pop_same
BenchmarkStack/UnsafeCompressedStack_Pop_same-16        1000000000               4.259 ns/op           0 B/op          0 allocs/op
BenchmarkStack/UnsafeUint32Stack_Push
BenchmarkStack/UnsafeUint32Stack_Push-16                45320574               147.7 ns/op            16 B/op          1 allocs/op
BenchmarkStack/UnsafeUint32Stack_Top
BenchmarkStack/UnsafeUint32Stack_Top-16                 1000000000               1.635 ns/op           0 B/op          0 allocs/op
BenchmarkStack/UnsafeUint32Stack_Pop
BenchmarkStack/UnsafeUint32Stack_Pop-16                 1000000000               2.507 ns/op           0 B/op          0 allocs/op
BenchmarkStack/UnsafeUint32Stack_Push_same
BenchmarkStack/UnsafeUint32Stack_Push_same-16           1000000000               3.093 ns/op           0 B/op          0 allocs/op
BenchmarkStack/UnsafeUint32Stack_Top_same
BenchmarkStack/UnsafeUint32Stack_Top_same-16            1000000000               1.977 ns/op           0 B/op          0 allocs/op
BenchmarkStack/UnsafeUint32Stack_Pop_same
BenchmarkStack/UnsafeUint32Stack_Pop_same-16            1000000000               2.706 ns/op           0 B/op          0 allocs/op
PASS
ok      github.com/lvyahui8/ellyn/sdk/common/collections        174.215s
```