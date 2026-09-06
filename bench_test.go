package riffle_test

import (
	"fmt"
	"testing"

	"github.com/Azridum/riffle"
)

var sizes = []int{10, 1_000, 100_000}
var sink int

func input(n int) []int {
	xs := make([]int, n)
	for i := range xs {
		xs[i] = i
	}
	return xs
}

func double(i int) int   { return i * 2 }
func even(i int) bool    { return i%2 == 0 }
func add(i1, i2 int) int { return i1 + i2 }

var pair = [2]int{1, 2}

func two(int) []int { return pair[:] }

func BenchmarkMap(b *testing.B) {
	for _, n := range sizes {
		xs := input(n)

		b.Run(fmt.Sprintf("loop/%d", n), func(b *testing.B) {
			for b.Loop() {
				r := make([]int, len(xs))
				for i := range xs {
					r[i] = double(xs[i])
				}
			}
		})

		b.Run(fmt.Sprintf("slice/%d", n), func(b *testing.B) {
			for b.Loop() {
				riffle.From(xs).Map(double)
			}
		})

		b.Run(fmt.Sprintf("seq/%d", n), func(b *testing.B) {
			for b.Loop() {
				riffle.From(xs).Seq().Map(double).Collect()
			}
		})
	}
}

func BenchmarkFilter(b *testing.B) {
	for _, n := range sizes {
		xs := input(n)

		b.Run(fmt.Sprintf("loop/%d", n), func(b *testing.B) {
			for b.Loop() {
				r := make([]int, 0, len(xs))
				for _, x := range xs {
					if even(x) {
						r = append(r, x)
						_ = r
					}
				}
			}
		})

		b.Run(fmt.Sprintf("slice/%d", n), func(b *testing.B) {
			for b.Loop() {
				riffle.From(xs).Filter(even)
			}
		})

		b.Run(fmt.Sprintf("seq/%d", n), func(b *testing.B) {
			for b.Loop() {
				riffle.From(xs).Seq().Filter(even).Collect()
			}
		})
	}
}

func BenchmarkReduce(b *testing.B) {
	for _, n := range sizes {
		xs := input(n)
		b.Run(fmt.Sprintf("loop/%d", n), func(b *testing.B) {
			for b.Loop() {
				v := xs[0]
				for _, x := range xs[1:] {
					v = add(v, x)
				}
				sink = v
			}
		})
		b.Run(fmt.Sprintf("slice/%d", n), func(b *testing.B) {
			for b.Loop() {
				sink, _ = riffle.From(xs).Reduce(add)
			}
		})
		b.Run(fmt.Sprintf("seq/%d", n), func(b *testing.B) {
			for b.Loop() {
				sink, _ = riffle.From(xs).Seq().Reduce(add)
			}
		})
	}
}

func BenchmarkFlatMap(b *testing.B) {
	for _, n := range sizes {
		xs := input(n)
		b.Run(fmt.Sprintf("loop/%d", n), func(b *testing.B) {
			for b.Loop() {
				r := make([]int, 0, len(xs))
				for _, x := range xs {
					r = append(r, two(x)...)
					_ = r
				}
			}
		})
		b.Run(fmt.Sprintf("slice/%d", n), func(b *testing.B) {
			for b.Loop() {
				riffle.From(xs).FlatMap(two)
			}
		})
		b.Run(fmt.Sprintf("seq/%d", n), func(b *testing.B) {
			for b.Loop() {
				riffle.From(xs).Seq().FlatMap(two).Collect()
			}
		})
	}
}
