package riffle_test

import (
	"testing"

	"github.com/Azridum/riffle"
	"github.com/shoenig/test"
	"pgregory.net/rapid"
)

func TestMapMatchesSlice(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		xs := rapid.SliceOf(rapid.Int()).Draw(t, "xs")
		f := func(i int) int { return i * 2 }

		eager := riffle.From(xs).Map(f)
		lazy := riffle.From(xs).Seq().Map(f).Collect()

		test.Eq(t, eager, lazy)

	})
}

func TestSeqMapDoesNotCallFnUntilConsumed(t *testing.T) {
	calls := 0
	riffle.Of(1, 2, 3).Seq().Map(func(i int) int {
		calls++
		return i
	})

	test.Eq(t, 0, calls)
}

func TestSeqMapStopsWhenConsumerStops(t *testing.T) {
	calls := 0
	for range riffle.Of(1, 2, 3).Seq().Map(func(i int) int {
		calls++
		return i
	}) {
		break
	}

	test.Eq(t, 1, calls)
}
