# Ellyn

Go coverage, call chain, and runtime data collection toolkit.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/lvyahui8/ellyn)](https://goreportcard.com/report/github.com/lvyahui8/ellyn)
[![codecov](https://codecov.io/gh/lvyahui8/ellyn/graph/badge.svg?token=YBV3TH2HQU)](https://codecov.io/gh/lvyahui8/ellyn)

**Languages:** English | [简体中文](README.zh-CN.md)

Ellyn instruments Go applications so you can collect request-level coverage,
function call chains, async execution links, and runtime data such as arguments,
return values, errors, and latency. It is designed for observability, precision
testing, traffic replay, and risk analysis scenarios where coverage alone is not
enough.

## Features

- Global coverage collection for full and incremental coverage analysis.
- Function call chain collection, including asynchronous links.
- Runtime data collection for parameters, return values, exceptions, and latency.
- Request-level collection for coverage, call graph, and runtime details.
- Concurrent data collection.
- Performance-oriented SDK components for high-frequency runtime paths.
- Experimental mock support is planned but not currently available.

## Use Cases

- Coverage statistics and single-test coverage details.
- Call chain tracing and runtime observation.
- Data and field lineage analysis.
- Traffic observation and replay.
- Precision testing and automated testing.
- Risk analysis.
- Unified metrics, monitoring, and alerting.

## Requirements

- Go 1.18 or later.
- Linux, macOS, or Windows.

## Demo

Download the demo program for your platform from
[GitHub Releases](https://github.com/lvyahui8/ellyn/releases), run it, and open
[http://localhost:19898](http://localhost:19898).

![Call chain visualization](./.assets/graph.png)

## CLI Usage

Download the `ellyn` CLI from
[GitHub Releases](https://github.com/lvyahui8/ellyn/releases).

```text
NAME:
   ellyn - Go coverage and callgraph collection tool

USAGE:
   ellyn [global options] command [command options]

COMMANDS:
   update    update code
   rollback  rollback code
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

Run the CLI in the directory that contains the target Go application's `main`
package:

```shell
ellyn update
```

`ellyn update` instruments the project. After instrumentation, compile and start
the target service to collect data.

```shell
ellyn rollback
```

`ellyn rollback` restores the original source files and removes instrumentation
artifacts.

## Repository Layout

- `api`: Runtime API used to access the instrumented SDK.
- `benchmark`: Benchmark scenarios for comparing overhead at different sampling rates.
- `cmd`: The `ellyn` command-line tool.
- `example`: Demo application for inspecting collected data.
- `instr`: Instrumentation logic that walks target Go files and injects SDK calls.
- `sdk`: SDK code copied into the target project and compiled as part of it.
- `test`: Shared test helpers.
- `viewer`: Lightweight visualization UI.

## Architecture

![Architecture](.assets/arch.png)

## SDK Flow

![SDK flow](.assets/flow.png)

## Development Notes

The SDK is copied into target projects, so runtime code should avoid unnecessary
dependencies and keep hot paths predictable:

- Prefer lock-free designs where practical, and avoid resource contention.
- Keep core operations at `O(1)`.
- Pad highly accessed fields to reduce false sharing.
- Trade limited memory overhead for lower runtime latency when it is justified.
- Prefer arrays and bitmaps over Go maps on high-frequency paths.
- Reuse frequently allocated objects with `sync.Pool` to reduce GC pressure.
- Be careful when collecting large parameter values because copying can be costly.

## SDK Components

- [RingBuffer](./sdk/common/collections/ringbuffer.go): Buffers call data.
  - [RingBuffer benchmark](./sdk/common/collections/ringbuffer.md)
  - [RingBuffer vs Map benchmark](./sdk/common/collections/ring_buffer_vs_map.md)
- [LinkedQueue](./sdk/common/collections/linked_queue.go): Linked-list-based synchronized queue used by the goroutine pool.
- [hmap / SegmentHashmap](./sdk/common/collections/hmap.go): High-performance routine-local storage implementation.
  - [hmap benchmark](./sdk/common/collections/hmap.md)
- [bitmap](./sdk/common/collections/bitmap.go): Records function and block execution.
- [UnsafeCompressedStack](./sdk/common/collections/stack.go): Simulates stack push and pop operations.
  - [Stack benchmark](./sdk/common/collections/stack.md)
- [routineLocal / GLS / GoRoutineLocalStorage](./sdk/common/goroutine/routine_local.go): Caches goroutine context.
  - [routineLocal benchmark](./sdk/common/goroutine/routine_local_test.go)
- [routinePool](./sdk/common/goroutine/routine_pool.go): Goroutine pool for concurrent file processing.
- [Uint64GUIDGenerator](./sdk/common/guid/guid.go): Generates call IDs.
- [AsyncLogger](./sdk/common/logging/readme.md): High-performance asynchronous logger.

## Performance

- CPU-intensive workloads can show measurable overhead even at low sampling rates.
- IO-intensive workloads are usually affected much less, even with higher sampling rates.

See [benchmark details](./benchmark/result.md).

## FAQ

### Why implement custom collection utilities instead of using open-source alternatives?

Part of the SDK is copied into target repositories. Custom implementations help
avoid dependency conflicts in those repositories and allow Ellyn to optimize for
its specific runtime collection paths. For code that is not copied into target
projects, reusing mature open-source implementations is preferred.
