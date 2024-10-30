package gocontext

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestContextTypHash(t *testing.T) {
	t.Log(valueCtxTypeHash)
	t.Log(emptyCtxTypeHash)
	t.Log(cancelCtxTypeHash)
	t.Log(timerCtxTypeHash)
}

func TestGetContextKeyValues(t *testing.T) {
	ctx := context.WithValue(context.Background(), 1, 2)
	ctx = context.WithValue(ctx, "name", "fg")
	ctx, cancelFunc := context.WithTimeout(ctx, time.Duration(3)*time.Second)
	cancelFunc()
	<-ctx.Done()
	ctx = context.WithValue(ctx, "age", 1)
	res := GetContextKeyValues(ctx)
	require.Equal(t, 3, len(res))
}
