# ellyn

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/lvyahui8/ellyn)](https://goreportcard.com/report/github.com/lvyahui8/ellyn)
[![codecov](https://codecov.io/gh/lvyahui8/ellyn/graph/badge.svg?token=YBV3TH2HQU)](https://codecov.io/gh/lvyahui8/ellyn)


### Requires

- Go Version >= 1.18

### key

- 避免资源冲突/锁竞争， 无锁优先
- 核心函数必须O(1)操作
- 高频访问的元素必须缓存行填充
- 牺牲牺牲部分空间换时间
- etc

### Sdk组件及用途

- [RingBuffer ](./ellyn_common/collections/ringbuffer.go) : 缓冲流量数据
  - [RingBuffer性能测试](./ellyn_common/collections/ringbuffer.md)
- [LinkedQueue](./ellyn_common/collections/linked_queue.go): 基于链表的同步队列。用作协程池的任务队列
- [hmap(SegmentHashmap)](./ellyn_common/collections/hmap.go): 实现高性能的routineLocal
  - [hmap性能测试](./ellyn_common/collections/hmap.md)
- [bitmap](./ellyn_common/collections/bitmap.go): 记录函数、块的执行情况
- [UnsafeCompressedStack](./ellyn_common/collections/stack.go) : 模拟入栈弹栈
  - [Stack性能测试](./ellyn_common/collections/stack.md)
- [routineLocal/GLS/GoRoutineLocalStorage](./ellyn_common/goroutine/routine_local.go): 缓存上下文
  - [routineLocal性能测试](./ellyn_common/gls/routine_local_test.go)
- [routinePool](./ellyn_common/goroutine/routine_pool.go): 协程池，并发处理文件
- [Uint64GUIDGenerator](./ellyn_common/guid/guid.go)

### Q&A 

#### Q:为什么要实现部分集合库，而不是直接使用开源方案？

A: 这部分实现会拷贝的目标仓库，为确保不与目标仓库sdk冲突，所以自行实现。对于不拷贝的目标仓库的代码，优先考虑用复用开源实现
