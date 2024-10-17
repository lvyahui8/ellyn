package ellyn_agent

import (
	"github.com/stretchr/testify/require"
	"testing"
	"unsafe"
)

func TestNotCollected(t *testing.T) {
	var a any = struct{}{}
	require.False(t, NotCollected == a)
	a = &struct{}{}
	require.False(t, NotCollected == a)
	a = struct {
		_ byte
	}{}
	require.False(t, NotCollected == a)
	a = &struct {
		_ byte
	}{}
	require.False(t, NotCollected == a)
	t.Log(unsafe.Sizeof(NotCollected))
	t.Log(unsafe.Sizeof(*NotCollected))
}
