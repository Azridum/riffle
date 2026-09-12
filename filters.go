package riffle

import "cmp"

// NotZeroValue reports whether v is not the zero value of its type.
//
//	nonZero := riffle.Of(4, 0, 2).Filter(riffle.NotZeroValue)
//	// nonZero == [4 2]
func NotZeroValue[T comparable](v T) bool {
	var zero T

	return v != zero
}

// Not returns a predicate that negates the result of fn.
//
//	even := func(i int) bool { return i%2 == 0 }
//	odds := riffle.Of(4, 1, 2, 3).Filter(riffle.Not(even))
//	// odds == [1 3]
func Not[T any](fn func(T) bool) func(T) bool {
	return func(v T) bool {
		return !fn(v)
	}
}

// GreaterThan returns a predicate that reports whether a value is strictly greater than val.
//
//	gt := riffle.Of(1, 2, 3, 4).Filter(riffle.GreaterThan(2))
//	// gt == [3 4]
func GreaterThan[T cmp.Ordered](val T) func(T) bool {
	return func(v T) bool {
		return v > val
	}
}

// GreaterThanOrEqual returns a predicate that reports whether a value is greater than or equal to val.
//
//	gte := riffle.Of(1, 2, 3, 4).Filter(riffle.GreaterThanOrEqual(2))
//	// gte == [2 3 4]
func GreaterThanOrEqual[T cmp.Ordered](val T) func(T) bool {
	return func(v T) bool {
		return v >= val
	}
}
