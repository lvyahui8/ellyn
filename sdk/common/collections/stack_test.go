package collections

import (
	"github.com/stretchr/testify/require"
	"testing"
	"unsafe"
)

type testInt uint32

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

func testStack[T testInt | uint32](b *testing.B, stackName string, stack Stack[T]) {
	b.Run(stackName+"_Push", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			stack.Push(T(i))
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
			stack.Push(1)
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

func TestUnsafeUint32Stack(t *testing.T) {
	s := NewUnsafeUint32Stack()
	s.Push(1)
	s.Push(1)
	s.Push(2)
	s.Push(3)
	s.Push(3)
	s.Push(4)
	// now 1[2]->2->3[2]-4
	val, suc := s.Pop()
	require.True(t, suc)
	require.Equal(t, uint32(4), val)

	// now  1[2]->2->3[2]
	val, suc = s.Pop()
	require.True(t, suc)
	require.Equal(t, uint32(3), val)

	// now  1[2]->2->3[1]
	val, suc = s.Pop()
	require.True(t, suc)
	require.Equal(t, uint32(3), val)

	// now  1[2]->2
	val, suc = s.Pop()
	require.True(t, suc)
	require.Equal(t, uint32(2), val)

	// now  1[2]
	val, suc = s.Pop()
	require.True(t, suc)
	require.Equal(t, uint32(1), val)

	// now  1[1]
	val, suc = s.Pop()
	require.True(t, suc)
	require.Equal(t, uint32(1), val)

	// now is empty
	val, suc = s.Pop()
	require.False(t, suc)
	require.Equal(t, uint32(0), val)
	require.True(t, s.Empty())
}

func TestUint32NodeSize(t *testing.T) {
	t.Log(unsafe.Sizeof(uint32Node{}.extra))
}

//  go test -v -run ^$  -bench BenchmarkStack -benchtime=5s  -benchmem
func BenchmarkStack(b *testing.B) {
	testStack[testInt](b, "UnsafeStack", NewUnsafeStack[testInt]())
	testStack[testInt](b, "UnsafeCompressedStack", NewUnsafeCompressedStack[testInt]())
	testStack[uint32](b, "UnsafeUint32Stack", NewUnsafeUint32Stack())
}

// BenchmarkUnsafeCompressedStack_Push
// go test -v -run ^$  -bench BenchmarkUnsafeCompressedStack_Push -benchtime=5s -benchmem -memprofile memprofile.pprof -cpuprofile profile.pprof
// go tool pprof -http=":8081" memprofile.pprof
func BenchmarkUnsafeCompressedStack_Push(b *testing.B) {
	stack := NewUnsafeCompressedStack[testInt]()
	for i := 0; i < b.N; i++ {
		stack.Push(1)
	}
}

func TestSizeofElementsPtr(t *testing.T) {
	stack := NewUnsafeCompressedStack[testInt]()
	t.Log(unsafe.Sizeof(stack.elements))
	t.Log(unsafe.Sizeof(&stackElement[testInt]{}))
	t.Log(unsafe.Sizeof(*stack))
}
