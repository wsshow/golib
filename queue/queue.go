package queue

type Queue[T comparable] struct {
	items []T
}

func New[T comparable](values ...T) *Queue[T] {
	qu := &Queue[T]{
		items: make([]T, 0),
	}
	return qu
}

func (q *Queue[T]) Enqueue(items ...T) {
	q.items = append(q.items, items...)
}

func (q *Queue[T]) Dequeue() interface{} {
	if q.IsEmpty() {
		return nil
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item
}

func (q *Queue[T]) Peek() T {
	return q.items[0]
}

func (q *Queue[T]) Count() int {
	return len(q.items)
}

func (q *Queue[T]) Contains(item T) bool {
	for _, qItem := range q.items {
		if qItem == item {
			return true
		}
	}
	return false
}

func (q *Queue[T]) ToSlice() []T {
	return q.items
}

func (q *Queue[T]) IsEmpty() bool {
	return q.Count() == 0
}

func (q *Queue[T]) Clear() {
	q.items = nil
}

func (q *Queue[T]) ForEach(callbackFn func(T)) {
	for _, qItem := range q.items {
		callbackFn(qItem)
	}
}

func (q *Queue[T]) Map(callbackFn func(T) T) *Queue[T] {
	nq := New[T]()
	for _, qItem := range q.items {
		nq.Enqueue(callbackFn(qItem))
	}
	return nq
}
