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

func TestGreater(t *testing.T) {
	result := riffle.Of(1, 2, 3, 4).Filter(riffle.GreaterThan(2))

	test.Eq(t, []int{3, 4}, result)

}

func TestGreaterThanOrEqual(t *testing.T) {
	result := riffle.Of(1, 2, 3, 4).Filter(riffle.GreaterThanOrEqual(2))

	test.Eq(t, []int{2, 3, 4}, result)

}
