package agent

import (
	"testing"
	"unsafe"
)

func TestFrameSize(t *testing.T) {
	t.Log(unsafe.Sizeof(methodFrame{}))
}
