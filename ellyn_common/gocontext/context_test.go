package gocontext

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
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
	res := GetContextKeyValues(ctx)
	require.Equal(t, 2, len(res))
}
