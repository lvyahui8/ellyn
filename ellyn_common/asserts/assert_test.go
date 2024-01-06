package asserts

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAssertNil(t *testing.T) {
	defer func() {
		e := recover()
		ae, ok := e.(AssertError)
		require.True(t, ok)
		require.NotNil(t, ae.msg)
		t.Log(ae.fileLine)
		t.Log(ae.msg)
	}()
	IsNil("xxx")
}
