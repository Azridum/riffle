package riffle_test

import (
	"testing"

	"github.com/Azridum/riffle"
	"github.com/shoenig/test"
)

func TestNotZeroValue(t *testing.T) {
	result := riffle.Of(4, 0, 2).Filter(riffle.NotZeroValue)
	test.Eq(t, []int{4, 2}, result)
}

func TestNot(t *testing.T) {
	even := func(v int) bool {
		return v%2 == 0
	}

	result := riffle.Of(4, 1, 2, 3).Filter(riffle.Not(even))

	test.Eq(t, []int{1, 3}, result)
}
