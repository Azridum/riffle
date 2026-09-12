package riffle_test

import (
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/Azridum/riffle"
	"github.com/shoenig/test"
	"pgregory.net/rapid"
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

func TestSliceFirst(t *testing.T) {
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
			v, ok := riffle.From(c.data).First()
			test.Eq(t, c.expected, v)
			test.Eq(t, c.expectedOk, ok)
		})
	}
}

func TestSliceLast(t *testing.T) {
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
			expected:   13,
			expectedOk: true,
		},
		"last element is zero": {
			data:       []int{11, 12, 0},
			expected:   0,
			expectedOk: true,
		},
	}

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			v, ok := riffle.From(c.data).Last()
			test.Eq(t, c.expected, v)
			test.Eq(t, c.expectedOk, ok)
		})
	}
}

func TestSliceReduce(t *testing.T) {
	cases := map[string]struct {
		data       []int
		expected   int
		expectedOk bool
		seen       [][2]int
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
			data:       []int{1},
			expected:   1,
			expectedOk: true,
		},
		"one zero element": {
			data:       []int{0},
			expected:   0,
			expectedOk: true,
		},
		"multiple elements": {
			data:       []int{1, 2, 3},
			expected:   123,
			expectedOk: true,
			seen:       [][2]int{{1, 2}, {12, 3}},
		},
	}

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			var seen [][2]int
			result, ok := riffle.From(c.data).Reduce(func(i1, i2 int) int {
				seen = append(seen, [2]int{i1, i2})
				return i1*10 + i2
			})

			test.Eq(t, c.expectedOk, ok)
			test.Eq(t, c.expected, result)
			test.Eq(t, c.seen, seen)
		})
	}
}

func TestSliceFold(t *testing.T) {
	type call struct {
		acc string
		v   int
	}

	cases := map[string]struct {
		data     []int
		init     string
		expected string
		seen     []call
	}{
		"nil slice": {
			data:     nil,
			init:     "1",
			expected: "1",
		},
		"empty slice": {
			data:     []int{},
			init:     "1",
			expected: "1",
		},
		"one element empty init": {
			data:     []int{2},
			init:     "",
			expected: "2",
			seen: []call{
				{
					acc: "",
					v:   2,
				},
			},
		},
		"one element not empty init": {
			data:     []int{2},
			init:     "1",
			expected: "12",
			seen: []call{
				{
					acc: "1",
					v:   2,
				},
			},
		},
		"multiple elements empty init": {
			data:     []int{2, 3, 4},
			init:     "",
			expected: "234",
			seen: []call{
				{
					acc: "",
					v:   2,
				},
				{
					acc: "2",
					v:   3,
				},
				{
					acc: "23",
					v:   4,
				},
			},
		},
		"multiple elements not empty init": {
			data:     []int{2, 3, 4},
			init:     "1",
			expected: "1234",
			seen: []call{
				{
					acc: "1",
					v:   2,
				},
				{
					acc: "12",
					v:   3,
				},
				{
					acc: "123",
					v:   4,
				},
			},
		},
	}

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			var seen []call
			test.Eq(t, c.expected, riffle.From(c.data).Fold(c.init, func(acc string, v int) string {
				seen = append(seen, call{
					acc: acc, v: v,
				})
				return acc + strconv.FormatInt(int64(v), 10)
			}))
			test.Eq(t, c.seen, seen)
		})
	}
}

func TestSliceFlatMap(t *testing.T) {
	cases := map[string]struct {
		data         []int
		callResults  []riffle.Slice[string]
		expected     riffle.Slice[string]
		expectedSeen []int
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
			data: []int{1},
			callResults: []riffle.Slice[string]{
				{"1", "1"},
			},
			expected:     []string{"1", "1"},
			expectedSeen: []int{1},
		},
		"multiple elements": {
			data: []int{1, 2, 3},
			callResults: []riffle.Slice[string]{
				{"1", "1"},
				{"2", "2"},
				{"3", "3"},
			},
			expected:     []string{"1", "1", "2", "2", "3", "3"},
			expectedSeen: []int{1, 2, 3},
		},
		"multiple elements some results empty": {
			data: []int{1, 2, 3},
			callResults: []riffle.Slice[string]{
				{"1", "1"},
				nil,
				{},
			},
			expected:     []string{"1", "1"},
			expectedSeen: []int{1, 2, 3},
		},
		"all results empty": {
			data: []int{1, 2, 3},
			callResults: []riffle.Slice[string]{
				nil,
				nil,
				nil,
			},
			expected:     nil,
			expectedSeen: []int{1, 2, 3},
		},
	}

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			var seen []int
			test.Eq(t, c.expected, riffle.From(c.data).FlatMap(func(i int) riffle.Slice[string] {
				seen = append(seen, i)
				return c.callResults[len(seen)-1]
			}))
			test.Eq(t, c.expectedSeen, seen)
		})
	}
}

func TestSliceFlatMapAcceptsPlainSliceFn(t *testing.T) {
	r := riffle.Of("a b", "c").FlatMap(strings.Fields)
	test.Eq(t, riffle.Slice[string]{"a", "b", "c"}, r)
}

func TestSliceGroupBy(t *testing.T) {
	type data struct {
		id   int
		name string
	}

	cases := map[string]struct {
		data     []data
		expected map[int]riffle.Slice[data]
	}{
		"nil slice": {
			data:     nil,
			expected: nil,
		},
		"empty slice": {
			data:     []data{},
			expected: nil,
		},
		"one element": {
			data: []data{
				{
					id:   1,
					name: "asd",
				},
			},
			expected: map[int]riffle.Slice[data]{
				1: {
					{
						id:   1,
						name: "asd",
					},
				},
			},
		},
		"multiple elements same id": {
			data: []data{
				{
					id:   1,
					name: "asd",
				},
				{
					id:   1,
					name: "asdasd",
				},
			},
			expected: map[int]riffle.Slice[data]{
				1: {
					{
						id:   1,
						name: "asd",
					},
					{
						id:   1,
						name: "asdasd",
					},
				},
			},
		},
		"multiple elements different id": {
			data: []data{
				{
					id:   1,
					name: "asd",
				},
				{
					id:   2,
					name: "asdasd",
				},
			},
			expected: map[int]riffle.Slice[data]{
				1: {
					{
						id:   1,
						name: "asd",
					},
				},
				2: {
					{
						id:   2,
						name: "asdasd",
					},
				},
			},
		},
	}

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			r := riffle.From(c.data).GroupBy(func(d data) int {
				return d.id
			})

			test.Eq(t, c.expected, r)
		})
	}
}

func TestGroupByIsAPartition(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		xs := rapid.SliceOf(rapid.Int()).Draw(t, "xs")
		key := func(i int) int { return i % 3 }

		groups := riffle.From(xs).GroupBy(key)

		total := 0
		for k, g := range groups {
			total += len(g)
			for _, v := range g {
				test.Eq(t, k, key(v))
			}
		}
		test.Eq(t, len(xs), total)
	})
}

func TestDistinct(t *testing.T) {
	cases := map[string]struct {
		data     []int
		expected riffle.Slice[int]
		seen     []int
	}{
		"nil slice": {
			data:     nil,
			expected: nil,
			seen:     nil,
		},
		"empty slice": {
			data:     []int{},
			expected: nil,
			seen:     nil,
		},
		"one element": {
			data:     []int{1},
			expected: []int{1},
			seen:     []int{1},
		},
		"multiple same elements": {
			data:     []int{1, 1},
			expected: []int{1},
			seen:     []int{1, 1},
		},
		"multiple different elements": {
			data:     []int{1, 2},
			expected: []int{1, 2},
			seen:     []int{1, 2},
		},
		"multiple not next to each other elements": {
			data:     []int{1, 2, 1, 2},
			expected: []int{1, 2},
			seen:     []int{1, 2, 1, 2},
		},
	}

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			var seen []int
			test.Eq(t, c.expected, riffle.From(c.data).Distinct(func(i int) int {
				seen = append(seen, i)
				return i
			}))
			test.Eq(t, c.seen, seen)
		})
	}
}

func TestFind(t *testing.T) {
	cases := map[string]struct {
		data       []int
		expected   int
		expectedOk bool
		seen       []int
	}{
		"nil slice": {
			data:       nil,
			expected:   0,
			expectedOk: false,
			seen:       nil,
		},
		"empty slice": {
			data:       []int{},
			expected:   0,
			expectedOk: false,
			seen:       nil,
		},
		"one element not found": {
			data:       []int{1},
			expected:   0,
			expectedOk: false,
			seen:       []int{1},
		},
		"one element found": {
			data:       []int{2},
			expected:   2,
			expectedOk: true,
			seen:       []int{2},
		},
		"multiple elements nothing found": {
			data:       []int{1, 1, 3},
			expected:   0,
			expectedOk: false,
			seen:       []int{1, 1, 3},
		},
		"multiple elements found": {
			data:       []int{1, 1, 4, 3},
			expected:   4,
			expectedOk: true,
			seen:       []int{1, 1, 4},
		},
	}

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			var seen []int

			v, ok := riffle.From(c.data).Find(func(i int) bool {
				seen = append(seen, i)
				return i%2 == 0
			})

			test.Eq(t, c.expectedOk, ok)
			test.Eq(t, c.expected, v)

			test.Eq(t, c.seen, seen)
		})
	}
}

func TestAdd(t *testing.T) {
	result, _ := riffle.Of(1, 0, 2, 3).Reduce(riffle.Sum)
	test.Eq(t, result, 6)
}
