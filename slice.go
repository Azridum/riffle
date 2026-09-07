// Package riffle provides chainable, generic collection operations over Go
// slices and iterators.
//
// Two types carry the API:
//
//   - [Slice] wraps a []T and evaluates every operation eagerly, allocating a
//     new Slice per step.
//   - [Seq] wraps an [iter.Seq] and evaluates lazily: building a pipeline does
//     no work, elements flow one at a time when the sequence is consumed, and
//     iteration stops as soon as the consumer stops.
//
// Both types expose the same operations and produce identical results; pick
// Slice when the data is small and you want the result as a slice, pick Seq
// for long chains, large inputs, or when you only need the first few results.
// Convert with [Slice.Seq] and [Seq.Collect].
//
//	sum, ok := riffle.Of(1, 2, 3, 4).
//		Seq().
//		Filter(func(i int) bool { return i%2 == 0 }).
//		Map(func(i int) int { return i * i }).
//		Reduce(func(a, b int) int { return a + b })
//	// sum == 20, ok == true
//
// Operations that may have nothing to return (First, Last, Reduce) report a
// second bool result instead of panicking or returning a pointer.
package riffle

import "slices"

// Slice is an eagerly evaluated []T. Every operation walks the whole slice
// and returns a new Slice; the receiver is never modified.
//
// A Slice is a plain slice type, so []T and Slice[T] convert freely and all
// the usual indexing, slicing, len and range work on it.
type Slice[T any] []T

// From wraps an existing slice without copying. Mutating data afterwards is
// visible through the returned Slice.
func From[T any](data []T) Slice[T] {
	return data
}

// Of builds a Slice from its arguments.
func Of[T any](data ...T) Slice[T] {
	return data
}

// Seq returns a lazy view over s. Elements are read from s at iteration time,
// so later changes to s are observed.
func (s Slice[T]) Seq() Seq[T] {
	return Seq[T](slices.Values(s))
}

// Map applies fn to every element, in order, and returns the results.
// Returns nil when s is empty.
func (s Slice[T]) Map[R any](fn func(T) R) Slice[R] {
	if len(s) == 0 {
		return nil
	}

	r := make(Slice[R], len(s))

	for i := range s {
		r[i] = fn(s[i])
	}

	return r
}

// Filter returns the elements for which fn reports true, preserving order.
// Returns nil when no element passes.
func (s Slice[T]) Filter(fn func(T) bool) Slice[T] {
	r := make(Slice[T], 0, len(s))

	for i := range s {
		if fn(s[i]) {
			r = append(r, s[i])
		}
	}

	if len(r) == 0 {
		return nil
	}

	return r
}

// First returns the first element and true, or the zero value and false when
// s is empty.
func (s Slice[T]) First() (T, bool) {
	if len(s) == 0 {
		var v T
		return v, false
	}

	return s[0], true
}

// Last returns the last element and true, or the zero value and false when
// s is empty.
func (s Slice[T]) Last() (T, bool) {
	if len(s) == 0 {
		var v T
		return v, false
	}

	return s[len(s)-1], true
}

// Reduce combines the elements left to right using fn, seeding the
// accumulator with the first element. Returns the zero value and false when
// s is empty; a single-element slice returns that element without calling fn.
//
// Use [Slice.Fold] when the accumulator has a different type than the
// elements or when an explicit initial value is needed.
func (s Slice[T]) Reduce(fn func(T, T) T) (T, bool) {
	if len(s) == 0 {
		var zero T
		return zero, false
	}
	v := s[0]
	values := s[1:]
	for i := range values {
		v = fn(v, values[i])
	}

	return v, true
}

// Fold combines the elements left to right into an accumulator of type R,
// starting from init. Returns init unchanged when s is empty.
func (s Slice[T]) Fold[R any](init R, fn func(R, T) R) R {
	for i := range s {
		init = fn(init, s[i])
	}

	return init
}

// FlatMap applies fn to every element and concatenates the resulting slices
// in order. fn may return any slice type whose underlying type is []R.
// Returns nil when the concatenation is empty.
func (s Slice[T]) FlatMap[R any, S ~[]R](fn func(T) S) Slice[R] {
	r := make(Slice[R], 0, len(s))

	for i := range s {
		o := fn(s[i])
		r = append(r, o...)
	}

	if len(r) == 0 {
		return nil
	}

	return r
}

func (s Slice[T]) GroupBy[K comparable](fn func(T) K) map[K]Slice[T] {
	if len(s) == 0 {
		return nil
	}

	r := make(map[K]Slice[T])

	for i := range s {
		k := fn(s[i])
		r[k] = append(r[k], s[i])
	}

	return r
}
