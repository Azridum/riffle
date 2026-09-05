package riffle_test

import (
	"testing"

	"github.com/Azridum/riffle"
	"github.com/shoenig/test"
)

func TestSliceMap(t *testing.T) {

	cases := map[string]struct {
		data     []int
		expected riffle.Slice[float64]
		fn       func(int) float64
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
