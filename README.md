# riffle

[![test](https://github.com/Azridum/riffle/actions/workflows/test.yml/badge.svg)](https://github.com/Azridum/riffle/actions/workflows/test.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/Azridum/riffle.svg)](https://pkg.go.dev/github.com/Azridum/riffle)
[![release](https://img.shields.io/github/v/release/Azridum/riffle)](https://github.com/Azridum/riffle/releases)

Chainable, generic `Map` / `Filter` / `Reduce` / `Fold` / `FlatMap` / `GroupBy` for Go slices and iterators. Two flavours with the same API:

| Type | Evaluation | Use when |
|------|-----------|----------|
| `Slice[T]` | eager, one allocation per step | small data, want a slice back |
| `Seq[T]` | lazy, wraps `iter.Seq[T]` | long chains, big inputs, early exit |

Requires Go 1.27+ (uses generic methods).

```sh
go get github.com/Azridum/riffle
```

## Usage

```go
import "github.com/Azridum/riffle"

// Eager
evens := riffle.Of(1, 2, 3, 4, 5, 6).
    Filter(func(i int) bool { return i%2 == 0 }).
    Map(func(i int) string { return strconv.Itoa(i) })
// riffle.Slice[string]{"2", "4", "6"}

// Lazy: nothing runs until First() pulls; stops after the 2nd element
first, ok := riffle.From(bigSlice).
    Seq().
    Map(expensive).
    Filter(func(x int) bool { return x > 100 }).
    First()

// Group into a map; each group keeps input order
byLen := riffle.Of("go", "rust", "zig", "odin").
    GroupBy(func(s string) int { return len(s) })
// map[int]riffle.Slice[string]{2: {"go"}, 3: {"zig"}, 4: {"rust", "odin"}}

// Seq is an iter.Seq: range over it directly
for v := range riffle.From(xs).Seq().Filter(pred) {
    fmt.Println(v)
}

// Convert between them
s.Seq()       // Slice -> Seq (no copy, reads s at iteration time)
seq.Collect() // Seq -> Slice
```

## API

Available on both `Slice[T]` and `Seq[T]` unless noted.

| Method | Description |
|--------|-------------|
| `Map(func(T) R) ...[R]` | transform each element |
| `Filter(func(T) bool)` | keep elements where fn is true |
| `FlatMap(func(T) S) ...[R]` with `S ~[]R` | map to slices, concatenate |
| `FlatMapSeq(func(T) S)` with `S ~func(func(R) bool)` | `Seq` only. Map to iterators, flatten lazily |
| `First() (T, bool)` | first element. `Seq` pulls exactly one |
| `Last() (T, bool)` | `Slice` only |
| `Reduce(func(T, T) T) (T, bool)` | fold seeded with first element; false when empty |
| `Fold(init R, func(R, T) R) R` | fold with explicit seed and accumulator type |
| `GroupBy(func(T) K) map[K]Slice[T]` with `K comparable` | partition by key; groups keep input order |
| `Distinct(func(T) K)` with `K comparable` | drop elements whose key was already seen; first occurrence wins, order kept |
| `Seq()` | `Slice` only. Lazy view |
| `Collect() Slice[T]` | `Seq` only. Drain to slice |

Constructors: `riffle.From([]T)` wraps without copying, `riffle.Of(a, b, c)` builds from arguments.

## Semantics worth knowing

- **Empty results are `nil`**, not `[]T{}`. `Map`, `Filter`, `FlatMap` and `Collect` all return `nil` when there is nothing to return. `GroupBy` returns a `nil` map on empty input; indexing it is fine, writing to it is not.
- **`Slice` and `Seq` always agree.** Every op is property-tested so that `s.Op(f)` equals `s.Seq().Op(f).Collect()`.
- **`Seq` is fully lazy.** Building a chain calls none of your functions. Consumers that stop early (`First`, `break` in a range loop) stop every upstream stage.
- **`Seq` is re-iterable.** Each consumption re-runs the pipeline from the source. Not safe for concurrent use unless the source is.
- **`From` does not copy.** Mutating the input afterwards is visible through the `Slice` and through any `Seq` built from it.
- **`Reduce` vs `Fold`.** `Reduce` seeds from the first element and reports `false` on empty input. `Fold` takes an explicit seed, can change type, and returns the seed on empty input.

## Development

```sh
go test ./...                            # unit + property tests (pgregory.net/rapid)
go test -bench . -run '^$' -benchmem     # loop vs slice vs seq benchmarks
```

Commit subjects follow [Conventional Commits](https://www.conventionalcommits.org/). A nightly workflow tags a new version from the commits since the last tag (`feat` bumps minor, `fix` bumps patch, `!` or a `BREAKING CHANGE` footer bumps major) and publishes the changelog on the [Releases](https://github.com/Azridum/riffle/releases) page. Other types (`docs`, `test`, `refactor`, `chore`, ...) never trigger a release on their own.

## License

MIT, see [LICENSE](LICENSE).
