package riffle_test

import (
	"slices"
	"testing"

	"github.com/Azridum/riffle"
	"github.com/shoenig/test"
)

func TestSliceMap(t *testing.T) {
	cases := map[string]struct {
		data     []int
		expected riffle.Slice[float64]
	}{
		"nil slice": {
			data:     nil,
			expected: nil,
		},
		"empty slice": {
			data:     []int{},
			expected: nil,
		},
		"one element": {
			data:     []int{1},
			expected: []float64{.5},
		},
		"multiple elements": {
			data:     []int{1, 2, 4},
			expected: []float64{.5, 1, 2},
		},
	}

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			s := riffle.From(c.data)

			r := s.Map(func(i int) float64 {
				return float64(i) / 2
			})

			test.Eq(t, c.expected, r)
		})
	}
}

func TestSliceMapCallsFnOncePerElementInOrder(t *testing.T) {
	data := []int{6, 1, 44}

	var seen []int
	r := riffle.From(data).Map(func(i int) int {
		seen = append(seen, i)
		return i * 10
	})

	test.Eq(t, data, seen)
	test.Eq(t, riffle.Slice[int]{60, 10, 440}, r)
}

func TestSliceMapDoesNotCallFnOnEmpty(t *testing.T) {
	calls := 0
	riffle.From([]int{}).Map(func(int) int {
		calls++
		return 0
	})

	test.Eq(t, 0, calls)
}

func TestSliceFilter(t *testing.T) {
	keep := func(int) bool {
		return true
	}
	drop := func(int) bool {
		return false
	}

	cases := map[string]struct {
		data     []int
		expected riffle.Slice[int]
		fn       func(int) bool
	}{
		"nil slice": {
			data:     nil,
			expected: nil,
			fn:       keep,
		},
		"empty slice": {
			data:     []int{},
			expected: nil,
			fn:       keep,
		},
		"one element nothing filtered": {
			data:     []int{1},
			expected: []int{1},
			fn:       keep,
		},
		"one element nothing left": {
			data:     []int{1},
			expected: nil,
			fn:       drop,
		},
		"multiple elements nothing filtered": {
			data:     []int{1, 2, 3},
			expected: []int{1, 2, 3},
			fn:       keep,
		},
		"multiple elements nothing left": {
			data:     []int{1, 2, 3},
			expected: nil,
			fn:       drop,
		},
		"multiple elements one filtered": {
			data:     []int{1, 2, 3},
			expected: []int{1, 3},
			fn: func(i int) bool {
				return i != 2
			},
		},
	}

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			old := slices.Clone(c.data)
			s := riffle.From(c.data)

			r := s.Filter(c.fn)

			test.Eq(t, c.expected, r)
			for i := range r {
				r[i] = -1
			}
			test.Eq(t, old, c.data)
		})
	}
}

func TestSliceFilterCallsFnOncePerElementInOrder(t *testing.T) {
	data := []int{6, 1, 44}

	var seen []int
	r := riffle.From(data).Filter(func(i int) bool {
		seen = append(seen, i)
		return true
	})

	test.Eq(t, data, seen)
	test.Eq(t, riffle.Slice[int]{6, 1, 44}, r)
}

func TestSliceFilterDoesNotCallFnOnEmpty(t *testing.T) {
	calls := 0
	riffle.From([]int{}).Filter(func(int) bool {
		calls++
		return true
	})

	test.Eq(t, 0, calls)
}
