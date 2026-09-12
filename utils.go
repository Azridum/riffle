package riffle

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
