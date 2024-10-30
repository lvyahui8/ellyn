package gocontext

import (
	"context"
	"unsafe"
)

func init() {
	emptyCtxTypeHash = contextTypeHash(context.Background())
	valueCtxTypeHash = contextTypeHash(context.WithValue(context.Background(), "", ""))
	ctx, cancelFunc := context.WithCancel(context.Background())
	go func() {
		cancelFunc()
	}()
	<-ctx.Done()
	cancelCtxTypeHash = contextTypeHash(ctx)
	tCtx, cancelFunc := context.WithTimeout(context.Background(), 0)
	cancelFunc()
	timerCtxTypeHash = contextTypeHash(tCtx)
}

type iface struct {
	tab  *itab
	data uintptr
}

type itab struct {
	inter uintptr
	_type *_type
}

type _type struct {
	_    uintptr
	_    uintptr // size of memory prefix holding all pointers
	hash uint32
}

type valueCtx struct {
	context.Context
	key, val interface{}
}

type cancelCtx struct {
	context.Context
}

type timerCtx struct {
	cancelCtx
}

func contextTypeHash(ctx context.Context) uint32 {
	f := (*iface)(unsafe.Pointer(&ctx))
	return f.tab._type.hash
}

var (
	valueCtxTypeHash  uint32
	emptyCtxTypeHash  uint32
	cancelCtxTypeHash uint32
	timerCtxTypeHash  uint32
)

func GetContextKeyValues(ctx context.Context) map[any]any {
	res := make(map[any]any)
	getKeyValue(ctx, res)
	return res
}

func getKeyValue(ctx context.Context, res map[any]any) {
	if ctx == nil {
		return
	}
	iCtx := (*iface)(unsafe.Pointer(&ctx))
	if iCtx.data == 0 {
		return
	}
	typeHash := iCtx.tab._type.hash
	switch typeHash {
	case valueCtxTypeHash:
		valCtx := (*valueCtx)(unsafe.Pointer(iCtx.data))
		if valCtx.Context == nil {
			return
		}
		if valCtx != nil && valCtx.key != nil && valCtx.val != nil {
			res[valCtx.key] = valCtx.val
		}
		getKeyValue(valCtx.Context, res)
	case emptyCtxTypeHash:
	case cancelCtxTypeHash:
		ccCtx := (*cancelCtx)(unsafe.Pointer(iCtx.data))
		getKeyValue(ccCtx.Context, res)
	case timerCtxTypeHash:
		tCtx := (*timerCtx)(unsafe.Pointer(iCtx.data))
		getKeyValue(tCtx.Context, res)
	}

}
