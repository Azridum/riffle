package riffle

import (
	"iter"
	"slices"
)

type Seq[T any] iter.Seq[T]

func (s Seq[T]) Collect() Slice[T] {
	return slices.Collect(iter.Seq[T](s))
}

func (s Seq[T]) Map[R any](fn func(T) R) Seq[R] {
	return func(yield func(R) bool) {
		for v := range s {
			if !yield(fn(v)) {
				return
			}
		}
	}
}

func (s Seq[T]) Filter(fn func(T) bool) Seq[T] {
	return func(yield func(T) bool) {
		for v := range s {
			if fn(v) && !yield(v) {
				return
			}
		}
	}
}
