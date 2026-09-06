package riffle_test

import (
	"testing"

	"github.com/Azridum/riffle"
	"github.com/shoenig/test"
	"pgregory.net/rapid"
)

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

func TestMapMatchesSlice(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		xs := rapid.SliceOf(rapid.Int()).Draw(t, "xs")
		f := func(i int) int { return i * 2 }

		eager := riffle.From(xs).Map(f)
		lazy := riffle.From(xs).Seq().Map(f).Collect()

		test.Eq(t, eager, lazy)
	})
}

func TestSeqFilterDoesNotCallFnUntilConsumed(t *testing.T) {
	calls := 0
	riffle.Of(1, 2, 3).Seq().Filter(func(i int) bool {
		calls++
		return true
	})

	test.Eq(t, 0, calls)
}

func TestSeqFilterStopsWhenConsumerStops(t *testing.T) {
	calls := 0
	for range riffle.Of(1, 2, 3).Seq().Filter(func(i int) bool {
		calls++
		return true
	}) {
		break
	}

	test.Eq(t, 1, calls)
}

func TestFilterMatchesSlice(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		xs := rapid.SliceOf(rapid.Int()).Draw(t, "xs")
		f := func(i int) bool { return i%2 == 0 }

		eager := riffle.From(xs).Filter(f)
		lazy := riffle.From(xs).Seq().Filter(f).Collect()

		test.Eq(t, eager, lazy)
	})
}
