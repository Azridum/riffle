package riffle

type Slice[T any] []T

func From[T any](data []T) Slice[T] {
	return data
}

func Of[T any](data ...T) Slice[T] {
	return data
}

func (s Slice[T]) Map[R any](fn func(T) R) Slice[R] {
	if len(s) == 0 {
		return nil
	}

	r := make(Slice[R], len(s))

	for i := range s {
		r[i] = fn(s[i])
	}

	return r
}

func (s Slice[T]) Filter(fn func(T) bool) Slice[T] {
	r := make(Slice[T], 0, len(s))

	for i := range s {
		if fn(s[i]) {
			r = append(r, s[i])
		}
	}

	if len(r) == 0 {
		return nil
	}

	return r
}

func (s Slice[T]) First() (T, bool) {
	if len(s) == 0 {
		var v T
		return v, false
	}

	return s[0], true
}
