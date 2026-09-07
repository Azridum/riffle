package riffle_test

import (
	"iter"
	"slices"
	"strconv"
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

func TestSeqFirst(t *testing.T) {
	cases := map[string]struct {
		data       []int
		expected   int
		expectedOk bool
	}{
		"nil slice": {
			data:       nil,
			expected:   0,
			expectedOk: false,
		},
		"empty slice": {
			data:       []int{},
			expected:   0,
			expectedOk: false,
		},
		"one element": {
			data:       []int{10},
			expected:   10,
			expectedOk: true,
		},
		"multiple elements": {
			data:       []int{11, 12, 13},
			expected:   11,
			expectedOk: true,
		},
		"first element is zero": {
			data:       []int{0, 12, 13},
			expected:   0,
			expectedOk: true,
		},
	}

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			v, ok := riffle.From(c.data).Seq().First()
			test.Eq(t, c.expected, v)
			test.Eq(t, c.expectedOk, ok)
		})
	}
}

func TestSeqFirstPullsOneElement(t *testing.T) {
	calls := 0
	v, ok := riffle.Of(1, 2, 3).Seq().Map(func(i int) int {
		calls++
		return i
	}).First()

	test.Eq(t, 1, calls)
	test.Eq(t, 1, v)
	test.True(t, ok)
}

func TestFirstMatchesSlice(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		xs := rapid.SliceOf(rapid.Int()).Draw(t, "xs")

		eager, eagerOk := riffle.From(xs).First()
		lazy, lazyOk := riffle.From(xs).Seq().First()

		test.Eq(t, eagerOk, lazyOk)
		test.Eq(t, eager, lazy)
	})
}

func TestReduceMatchesSlice(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		xs := rapid.SliceOf(rapid.Int()).Draw(t, "xs")

		fn := func(i1, i2 int) int {
			return i1*10 + i2
		}
		eager, eagerOk := riffle.From(xs).Reduce(fn)
		lazy, lazyOk := riffle.From(xs).Seq().Reduce(fn)

		test.Eq(t, eagerOk, lazyOk)
		test.Eq(t, eager, lazy)
	})
}

func TestFoldMatchesSlice(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		xs := rapid.SliceOf(rapid.Int()).Draw(t, "xs")

		fn := func(acc string, v int) string {
			return acc + strconv.FormatInt(int64(v), 10)
		}
		eager := riffle.From(xs).Fold("asd", fn)
		lazy := riffle.From(xs).Seq().Fold("asd", fn)

		test.Eq(t, eager, lazy)
	})
}

func TestFlatMapMatchesSlice(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		xs := rapid.SliceOf(rapid.Int()).Draw(t, "xs")

		fn := func(i int) []string {
			v := strconv.FormatInt(int64(i), 10)

			return []string{v, v + "asd"}
		}
		eager := riffle.From(xs).FlatMap(fn)
		lazy := riffle.From(xs).Seq().FlatMap(fn).Collect()
		lazySeq := riffle.From(xs).Seq().FlatMapSeq(func(i int) iter.Seq[string] {
			return slices.Values(fn(i))
		}).Collect()

		test.Eq(t, eager, lazy)
		test.Eq(t, eager, lazySeq)
	})
}

func TestSeqFlatMapDoesNotCallFnUntilConsumed(t *testing.T) {
	calls := 0
	riffle.Of(1, 2, 3).Seq().FlatMap(func(i int) []int {
		calls++
		return []int{i, i}
	})
	test.Eq(t, 0, calls)
}

func TestSeqFlatMapStopsWhenConsumerStops(t *testing.T) {
	var seen, got []int
	for v := range riffle.Of(1, 2, 3).Seq().FlatMap(func(i int) []int {
		seen = append(seen, i)
		return []int{i, i * 10}
	}) {
		got = append(got, v)
		if len(got) == 3 {
			break
		}
	}
	test.Eq(t, []int{1, 2}, seen)
	test.Eq(t, []int{1, 10, 2}, got)
}

func TestSeqFlatMapSeqStopsInnerWhenConsumerStops(t *testing.T) {
	var seen, got []int
	innerPulls := 0
	for v := range riffle.Of(1, 2, 3).Seq().FlatMapSeq(func(i int) riffle.Seq[int] {
		seen = append(seen, i)
		return func(yield func(int) bool) {
			for n := 0; ; n++ {
				innerPulls++
				if !yield(i*10 + n) {
					return
				}
			}
		}
	}) {
		got = append(got, v)
		if len(got) == 2 {
			break
		}
	}
	test.Eq(t, []int{1}, seen)
	test.Eq(t, 2, innerPulls)
	test.Eq(t, []int{10, 11}, got)
}

func TestGroupByMatchesSlice(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		xs := rapid.SliceOf(rapid.Int()).Draw(t, "xs")
		key := func(i int) int { return i % 3 }

		eager := riffle.From(xs).GroupBy(key)
		lazy := riffle.From(xs).Seq().GroupBy(key)

		test.Eq(t, eager, lazy)
	})
}
