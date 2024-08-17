package collections

type Queue interface {
	Enqueue(value any) (success bool)
	Dequeue() (value any, success bool)
}
