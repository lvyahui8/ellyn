package goroutine

import "runtime"

// GetGoId 获取协程id
// 依赖先执行拷贝动作，将EllynGetGoid方法写入go源码目录，可以参考脚本
// - goid_init.sh
// - goid_init.bat
func GetGoId() uint64 {
	return runtime.EllynGetGoid()
}

func GetRoutineCtx() uintptr {
	return runtime.EllynGetRoutineCtx()
}

func SetRoutineCtx(ctx uintptr) {
	runtime.EllynSetRoutineCtx(ctx)
}
