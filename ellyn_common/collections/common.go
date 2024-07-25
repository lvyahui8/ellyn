package collections

// https://www.lenshood.dev/2021/04/19/lock-free-ring-buffer/

func roundingToPowerOfTwo(size uint64) uint64 {
	size--
	size |= size >> 1
	size |= size >> 2
	size |= size >> 4
	size |= size >> 8
	size |= size >> 16
	size |= size >> 32
	size++
	return size
}
