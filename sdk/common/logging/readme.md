
## AsyncLogger

高性能异步日志，支持自动基于日期、大小切分日志文件，并支持自动清理

### Usage 

```go
initLogger()
log.Info("name:%s|age:%d|suc:%b", "yah", 1, true)
log.InfoKV(Empty().Str("name", "yah").Int("age", 1).Bool("suc", true))
log.InfoKV(logging.Code("g_collect").Int("n", len(g.nodes)).
Int("e", len(g.edges)).Bool("c", g.origin != nil))
```

### 关键性能优化

- Lock-Free RingBuffer实现的异步日志
- Buffer池化复用
- bufio.Writer缓冲写磁盘
- 单线程写文件，不加锁，减少性能消耗
- 缓存日期时间，减少系统调用或格式化日期耗时

### 与Go官方log库性能对比

#### 测试代码

```go
// go test -v -run ^$  -bench 'BenchmarkLogger' -benchtime=5s -benchmem
func BenchmarkLogger(b *testing.B) {
	initLogger()
	b.Run("asyncLogger", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			log.Info("hello world")
		}
	})
	b.Run("syncLogger", func(b *testing.B) {
		f, err := os.OpenFile("logs/test.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		asserts.IsNil(err)
		defer f.Close()
		ll.SetOutput(f)
		for i := 0; i < b.N; i++ {
			ll.Println("hello world")
		}
	})
}

func BenchmarkLoggerFormatAndKV(b *testing.B) {
	initLogger()
	b.Run("formatLog", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			log.Info("name:%s|age:%d|suc:%b", "yah", 1, true)
		}
	})
	b.Run("kvLog", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			log.InfoKV(Empty().Str("name", "yah").Int("age", 1).Bool("suc", true))
		}
	})
}
```

#### 测试结果

- asyncLogger写日志的耗时是go官方库写日志耗时的1.3%（性能相差约85倍）
- asyncLogger的kvLog比formatLog更快，耗时可以再降60%.

```text
goos: windows
goarch: amd64
pkg: github.com/lvyahui8/ellyn/sdk/common/logging
BenchmarkLogger
BenchmarkLogger/asyncLogger
BenchmarkLogger/asyncLogger-16          84107826                72.27 ns/op            0 B/op          0 allocs/op
BenchmarkLogger/syncLogger
BenchmarkLogger/syncLogger-16            1009458              5311 ns/op              32 B/op          1 allocs/op
BenchmarkLoggerFormatAndKV
BenchmarkLoggerFormatAndKV/formatLog
BenchmarkLoggerFormatAndKV/formatLog-16                 19870562               318.9 ns/op            54 B/op          1 allocs/op
BenchmarkLoggerFormatAndKV/kvLog
BenchmarkLoggerFormatAndKV/kvLog-16                     47905605               119.1 ns/op             0 B/op          0 allocs/op
PASS
ok      github.com/lvyahui8/ellyn/sdk/common/logging    28.942s

```
