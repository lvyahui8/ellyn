package collections

type Queue[T any] interface {
	Enqueue(value T) (success bool)
	Dequeue() (value T, success bool)
}
