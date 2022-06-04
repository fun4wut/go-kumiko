package yuriko

// Result implements https://doc.rust-lang.org/std/result/#structs
type Result[T any] struct {
	data T
	err  error
}

func Ok[T any](data T) Result[T] {
	return Result[T]{data: data}
}

func Err[T any](err error) Result[T] {
	return Result[T]{err: err}
}

func (r Result[T]) isError() bool {
	return r.err != nil
}

func (r *Result[T]) Unwrap() T {
	if r.isError() {
		panic("gg")
	}
	return r.data
}

func (r *Result[T]) UnwrapOr(def T) T {
	if r.isError() {
		return def
	}
	return r.data
}

func Map[T any, E any](r Result[T], f func(d T) E) Result[E] {
	if r.isError() {
		return Err[E](r.err)
	}
	return Ok[E](f(r.Unwrap()))
}

func AndThen[T any, E any](r Result[T], f func(d T) Result[E]) Result[E] {
	if r.isError() {
		return Err[E](r.err)
	}
	return f(r.Unwrap())
}
