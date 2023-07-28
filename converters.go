package croconf

func StringToCustomType[T any](source TypedBinding[string], callback func(string) (T, error)) TypedBinding[T] {
	return ToBinding[T](source.Identifier(), source.Source(), func() (T, error) {
		val, err := source.GetValue()
		if err != nil {
			var noop T
			return noop, err
		}
		return callback(val)
	})
}

func ByteSliceToCustomType[T any](source TypedBinding[[]byte], callback func([]byte) (T, error)) TypedBinding[T] {
	return ToBinding[T](source.Identifier(), source.Source(), func() (T, error) {
		val, err := source.GetValue()
		if err != nil {
			var noop T
			return noop, err
		}
		return callback(val)
	})
}
