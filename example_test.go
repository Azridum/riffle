package riffle_test

import (
	"fmt"
	"strings"

	"github.com/Azridum/riffle"
)

func ExampleSlice() {
	words := riffle.Of("go", "rust", "zig", "c").
		Filter(func(s string) bool { return len(s) > 1 }).
		Map(strings.ToUpper)

	fmt.Println(words)
	// Output: [GO RUST ZIG]
}

func ExampleSeq() {
	calls := 0
	first, ok := riffle.Of(1, 2, 3, 4).
		Seq().
		Map(func(i int) int { calls++; return i * 10 }).
		Filter(func(i int) bool { return i > 15 }).
		First()

	fmt.Println(first, ok, calls)
	// Output: 20 true 2
}

func ExampleSeq_Collect() {
	squares := riffle.Of(1, 2, 3).
		Seq().
		Map(func(i int) int { return i * i }).
		Collect()

	fmt.Println(squares, len(squares))
	// Output: [1 4 9] 3
}

func ExampleSlice_Reduce() {
	_, ok := riffle.Of[int]().Reduce(func(a, b int) int { return a + b })
	fmt.Println(ok)

	sum, ok := riffle.Of(1, 2, 3).Reduce(func(a, b int) int { return a + b })
	fmt.Println(sum, ok)
	// Output:
	// false
	// 6 true
}

func ExampleSlice_Fold() {
	csv := riffle.Of(1, 2, 3).Fold("", func(acc string, i int) string {
		if acc == "" {
			return fmt.Sprint(i)
		}
		return acc + "," + fmt.Sprint(i)
	})

	fmt.Println(csv)
	// Output: 1,2,3
}

func ExampleSeq_FlatMap() {
	chars := riffle.Of("ab", "c").
		Seq().
		FlatMap(func(s string) []string { return strings.Split(s, "") }).
		Collect()

	fmt.Println(chars)
	// Output: [a b c]
}

func ExampleSlice_GroupBy() {
	byLen := riffle.Of("go", "rust", "zig", "c", "odin").
		GroupBy(func(s string) int { return len(s) })

	fmt.Println(byLen)
	fmt.Println(byLen[4], len(byLen[7]))
	// Output:
	// map[1:[c] 2:[go] 3:[zig] 4:[rust odin]]
	// [rust odin] 0
}

func ExampleSeq_GroupBy() {
	parity := riffle.Of(1, 2, 3, 4, 5).
		Seq().
		Map(func(i int) int { return i * i }).
		GroupBy(func(i int) string {
			if i%2 == 0 {
				return "even"
			}
			return "odd"
		})

	fmt.Println(parity)
	// Output: map[even:[4 16] odd:[1 9 25]]
}

func ExampleSlice_Distinct() {
	byLen := riffle.Of("go", "rust", "zig", "c", "odin").
		Distinct(func(s string) int { return len(s) })

	fmt.Println(byLen)
	// Output: [go rust zig c]
}

func ExampleSeq_Distinct() {
	names := riffle.Of("Go", "go", "GO", "Rust").
		Seq().
		Distinct(strings.ToLower).
		Collect()

	fmt.Println(names)
	// Output: [Go Rust]
}
