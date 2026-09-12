package riffle

import (
	"cmp"
)

// Sum returns the sum of a and b.
//
//	sum, ok := riffle.Of(1, 0, 2, 3).Reduce(riffle.Sum)
//	// sum == 6, ok == true
func Sum[T number](a T, b T) T {
	return a + b
}

// Min returns the smaller of a and b.
//
//	lowest, ok := riffle.Of(4, 1, 2, 3).Reduce(riffle.Min)
//	// lowest == 1, ok == true
func Min[T cmp.Ordered](a T, b T) T {
	return min(a, b)
}

// Max returns the highest of a and b.
//
//	highest, ok := riffle.Of(4, 1, 2, 3).Reduce(riffle.Max)
//	// highest == 4, ok == true
func Max[T cmp.Ordered](a T, b T) T {
	return max(a, b)
}
