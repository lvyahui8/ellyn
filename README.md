# ellyn

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/lvyahui8/ellyn)](https://goreportcard.com/report/github.com/lvyahui8/ellyn)
[![codecov](https://codecov.io/gh/lvyahui8/ellyn/graph/badge.svg?token=YBV3TH2HQU)](https://codecov.io/gh/lvyahui8/ellyn)


### Requires

- Go Version >= 1.18


### Sdk组件及用途

- [RingBuffer ](./ellyn_common/collections/ringbuffer.go) : 缓冲流量数据
  - [RingBuffer性能测试](./ellyn_common/collections/ringbuffer.md)
- [hmap(SegmentHashmap)](./ellyn_common/collections/hmap.go): 实现高性能的routineLocal
  - [hmap性能测试](./ellyn_common/collections/hmap.md)
- [bitmap](./ellyn_common/collections/bitmap.go): 记录函数、块的执行情况
- [UnsafeCompressedStack](./ellyn_common/collections/stack.go) : 模拟入栈弹栈
  - [Stack性能测试](./ellyn_common/collections/stack.md)
- [routineLocal/GLS/GoRoutineLocalStorage](./ellyn_common/goroutine/routine_local.go): 缓存上下文
  - [routineLocal性能测试](./ellyn_common/gls/routine_local_test.go)
- [routinePool](./ellyn_common/goroutine/routine_pool.go): 协程池，并发处理文件
