package riffle_test

import (
	"testing"

	"github.com/Azridum/riffle"
	"github.com/shoenig/test"
)

func TestSum(t *testing.T) {
	result, _ := riffle.Of(1, 0, 2, 3).Reduce(riffle.Sum)
	test.Eq(t, result, 6)

	result = riffle.Of(1, 0, 2, 3).Fold(0.0, riffle.Sum)
	test.Eq(t, result, 6)
}

func TestMin(t *testing.T) {
	result, _ := riffle.Of(4, 1, 2, 3).Reduce(riffle.Min)
	test.Eq(t, result, 1)
}

func TestMax(t *testing.T) {
	result, _ := riffle.Of(4, 1, 2, 3).Reduce(riffle.Max)
	test.Eq(t, result, 4)
}
