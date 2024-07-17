package cond

func Some[T any](value T) Option[T] {
	return Option[T]{Some: true, Value: value}
}

func None[T any]() Option[T] {
	return Option[T]{Some: false}
}

// Option provides an optional value.
type Option[T any] struct {
	// Some is set to true when the option has a value.
	Some bool

	// Value is valid if Some is true.
	Value T
}

func (o Option[T]) MustGet(ifMissing T) T {
	if o.Some == true {
		return o.Value
	} else {
		return ifMissing
	}
}
