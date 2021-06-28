package croconf

type SourceCLI struct {
	flags []string
	// TODO
}

func NewSourceFromCLIFlags(flags []string) *SourceCLI {
	return &SourceCLI{flags: flags}
}

func (sc *SourceCLI) GetName() string {
	return "CLI flags" // TODO
}

func (sc *SourceCLI) FromName(name string) LazySingleValueBinding {
	return &cliBinding{source: sc, longhand: name}
}

func (sc *SourceCLI) FromNameAndShorthand(name, shorthand string) LazySingleValueBinding {
	return &cliBinding{source: sc, longhand: name, shorthand: shorthand}
}

func (sc *SourceCLI) FromPositionalArg(position int) LazySingleValueBinding {
	return &cliBinding{source: sc, position: position}
}

type cliBinding struct {
	source    *SourceCLI
	isUnary   bool // TODO: figure out what we should do about boolean CLI flags
	shorthand string
	longhand  string
	position  int
}

func (cb *cliBinding) GetSource() Source {
	return cb.source
}

func (cb *cliBinding) SaveStringTo(dest *string) error {
	return ErrorMissing // TODO: implement
}

func (cb *cliBinding) SaveInt64To(dest *int64) error {
	return ErrorMissing // TODO: implement
}

func (cb *cliBinding) SaveCustomValueTo(dest CustomValue) error {
	return ErrorMissing // TODO: implement
}
