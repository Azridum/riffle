package riffle

import (
	"iter"
	"slices"
)

// Seq is a lazily evaluated sequence of T, compatible with [iter.Seq].
//
// Building a pipeline with Map, Filter, FlatMap and friends performs no work
// and calls none of the supplied functions. Elements are produced one at a
// time when the Seq is consumed by a range loop or a terminal operation
// ([Seq.Collect], [Seq.First], [Seq.Reduce], [Seq.Fold], [Seq.GroupBy]). When
// the consumer stops early, every stage stops too, so no element beyond the
// ones actually needed is ever processed.
//
// A Seq can be iterated more than once; each iteration re-runs the pipeline
// from the source. Like iter.Seq, a Seq is not safe for concurrent use unless
// its source is.
type Seq[T any] iter.Seq[T]

// Collect drains s into a Slice. Returns nil when s yields no elements.
func (s Seq[T]) Collect() Slice[T] {
	return slices.Collect(iter.Seq[T](s))
}

// Map returns a Seq that yields fn applied to each element of s. fn is called
// lazily, once per element consumed.
func (s Seq[T]) Map[R any](fn func(T) R) Seq[R] {
	return func(yield func(R) bool) {
		for v := range s {
			if !yield(fn(v)) {
				return
			}
		}
	}
}

// Filter returns a Seq that yields only the elements of s for which fn
// reports true. fn is called lazily, once per element pulled from s.
func (s Seq[T]) Filter(fn func(T) bool) Seq[T] {
	return func(yield func(T) bool) {
		for v := range s {
			if fn(v) && !yield(v) {
				return
			}
		}
	}
}

// First consumes at most one element from s and returns it with true, or the
// zero value and false when s is empty. Upstream stages run only as far as
// needed to produce that element.
func (s Seq[T]) First() (T, bool) {
	for v := range s {
		return v, true
	}

	var zero T
	return zero, false
}

// Reduce drains s, combining the elements left to right using fn and seeding
// the accumulator with the first element. Returns the zero value and false
// when s is empty; a single-element sequence returns that element without
// calling fn.
//
// Use [Seq.Fold] when the accumulator has a different type than the elements
// or when an explicit initial value is needed.
func (s Seq[T]) Reduce(fn func(T, T) T) (T, bool) {
	var r T
	ok := false
	for v := range s {
		if !ok {
			r, ok = v, true
		} else {
			r = fn(r, v)
		}
	}

	return r, ok
}

// Fold drains s, combining the elements left to right into an accumulator of
// type R starting from init. Returns init unchanged when s is empty.
func (s Seq[T]) Fold[R any](init R, fn func(R, T) R) R {
	for v := range s {
		init = fn(init, v)
	}

	return init
}

// FlatMap returns a Seq that yields, for each element of s, every element of
// the slice returned by fn, in order. fn may return any slice type whose
// underlying type is []R. fn is called lazily, once per element pulled from s.
//
// Use [Seq.FlatMapSeq] when fn produces an iterator instead of a slice.
func (s Seq[T]) FlatMap[R any, S ~[]R](fn func(T) S) Seq[R] {
	return func(yield func(R) bool) {
		for v := range s {
			for _, r := range fn(v) {
				if !yield(r) {
					return
				}
			}
		}
	}
}

// FlatMapSeq is like [Seq.FlatMap] but fn returns an iterator, so the inner
// sequences are also produced lazily. fn may return an [iter.Seq], a [Seq],
// or any other func(yield func(R) bool).
func (s Seq[T]) FlatMapSeq[R any, S ~func(func(R) bool)](fn func(T) S) Seq[R] {
	return func(yield func(R) bool) {
		for v := range s {
			for r := range fn(v) {
				if !yield(r) {
					return
				}
			}
		}
	}
}

// GroupBy drains s and partitions the elements by the key fn returns for
// each of them. Every element appears in exactly one group, and each group
// keeps its elements in the order they were yielded. fn is called once per
// element, in order. Returns nil when s yields no elements.
//
// Each group is a freshly allocated Slice that shares no memory with the
// source of s.
func (s Seq[T]) GroupBy[K comparable](fn func(T) K) map[K]Slice[T] {
	r := make(map[K]Slice[T])

	for v := range s {
		k := fn(v)
		r[k] = append(r[k], v)
	}
	if len(r) == 0 {
		return nil
	}

	return r
}
