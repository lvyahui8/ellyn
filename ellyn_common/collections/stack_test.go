package collections

import "testing"

func testStack(b *testing.B, stackName string, stack Stack) {
	b.Run(stackName+"_Push", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			stack.Push(i)
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

func BenchmarkStack(b *testing.B) {
	testStack(b, "UnsafeStack", NewUnsafeStack())
	testStack(b, "UnsafeCompressedStack", NewUnsafeCompressedStack())
}
