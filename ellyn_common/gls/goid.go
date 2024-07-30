package gls

import "runtime"

func GetGoId() uint64 {
	return runtime.EllynGetGoid()
}
