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
