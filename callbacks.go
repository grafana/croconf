package croconf

type callbackBinding struct {
	apply func() error
}

func (cb *callbackBinding) Apply() error {
	return cb.apply()
}

func NewCallbackBinding(callback func() error) Binding {
	return &callbackBinding{apply: callback}
}

type callbackBindingFromSource struct {
	*callbackBinding
	source    Source
	boundName string
}

func (cbs *callbackBindingFromSource) Source() Source {
	return cbs.source
}

func (cbs *callbackBindingFromSource) BoundName() string {
	return cbs.boundName
}

func NewCallbackBindingFromSource(source Source, boundName string, callback func() error) BindingFromSource {
	return &callbackBindingFromSource{
		callbackBinding: &callbackBinding{apply: callback},
		source:          source,
		boundName:       boundName,
	}
}

func wrapBinding(origBinding Binding, newCallback func() error) Binding {
	if fromSource, ok := origBinding.(BindingFromSource); ok {
		return NewCallbackBindingFromSource(fromSource.Source(), fromSource.BoundName(), newCallback)
	} else {
		return NewCallbackBinding(newCallback)
	}
}
