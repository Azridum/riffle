package riffle

type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~complex64 | ~complex128
}

// Sum returns the sum of a and b.
//
//	sum, ok := riffle.Of(1, 0, 2, 3).Reduce(riffle.Sum)
//	// sum == 6, ok == true
func Sum[T number](a T, b T) T {
	return a + b
}
