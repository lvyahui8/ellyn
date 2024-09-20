package collections

import (
	"github.com/stretchr/testify/require"
	"testing"
	"unsafe"
)

type testInt int

func (t testInt) Equals(value Frame) bool {
	val, ok := value.(testInt)
	if ok {
		return t == val
	}
	return false
}

func (t testInt) Init() {
}

func (t testInt) ReEnter() {
}

func testStack(b *testing.B, stackName string, stack Stack[testInt]) {
	b.Run(stackName+"_Push", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			stack.Push(testInt(i))
		}
	})
	b.Run(stackName+"_Top", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			stack.Top()
		}
	})
	b.Run(stackName+"_Pop", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			stack.Pop()
		}
	})
	b.Run(stackName+"_Push_same", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			stack.Push(-1)
		}
	})
	b.Run(stackName+"_Top_same", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			stack.Top()
		}
	})
	b.Run(stackName+"_Pop_same", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			stack.Pop()
		}
	})
}

func TestUnsafeCompressedStack_Push(t *testing.T) {
	stack := NewUnsafeCompressedStack[testInt]()
	stack.Push(1)
	stack.Push(1)
	require.Equal(t, 2, stack.count)
	require.Equal(t, 1, stack.elements.Len())
}

//
func BenchmarkStack(b *testing.B) {
	testStack(b, "UnsafeStack", NewUnsafeStack[testInt]())
	testStack(b, "UnsafeCompressedStack", NewUnsafeCompressedStack[testInt]())
}

// BenchmarkUnsafeCompressedStack_Push
// go test -v -run ^$  -bench BenchmarkUnsafeCompressedStack_Push -benchtime=5s -benchmem -memprofile memprofile.pprof -cpuprofile profile.pprof
// go tool pprof -http=":8081" memprofile.pprof
func BenchmarkUnsafeCompressedStack_Push(b *testing.B) {
	stack := NewUnsafeCompressedStack[testInt]()
	for i := 0; i < b.N; i++ {
		stack.Push(-1)
	}
}

func TestSizeofElementsPtr(t *testing.T) {
	stack := NewUnsafeCompressedStack[testInt]()
	t.Log(unsafe.Sizeof(stack.elements))
	t.Log(unsafe.Sizeof(&stackElement[testInt]{}))
	t.Log(unsafe.Sizeof(*stack))
}
