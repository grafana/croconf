package croconf

type customBinding[T any] struct {
	id       string
	getValue func() (T, error)
	source   Source
}

func (cb *customBinding[T]) Identifier() string {
	return cb.id
}

func (cb *customBinding[T]) Source() Source {
	return cb.source
}

func (cb *customBinding[T]) GetValue() (T, error) {
	return cb.getValue()
}

func ToBinding[T any](identifier string, source Source, getValue func() (T, error)) TypedBinding[T] {
	return &customBinding[T]{
		id:       identifier,
		getValue: getValue,
		source:   source,
	}
}
