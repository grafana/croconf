package croconf

type Binding interface {
	Source() Source
	Identifier() string
}

type TypedBinding[T any] interface {
	Binding
	GetValue() (T, error)
}

type Source interface {
	Initialize() error // TODO: figure out a better 2-step mechanism?
	GetName() string   // TODO: remove?
}

// TODO: are these even necessary?
type Field interface {
	Destination() interface{}
	Name() string
	Bindings() []Binding
	Consolidate() (ConsolidatedField, error)
}

type ConsolidatedField interface {
	HasBeenSet() bool
	Source() Source
	Validate() error // TODO: make this part of consolidation?
}
