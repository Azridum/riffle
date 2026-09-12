package riffle

// NotZeroValue reports whether v is not the zero value of its type.
//
//	nonZero := riffle.Of(4, 0, 2).Filter(riffle.NotZeroValue)
//	// nonZero == [4 2]
func NotZeroValue[T comparable](v T) bool {
	var zero T

	return v != zero
}
