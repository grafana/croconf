package croconf

type SourceCLI struct {
	flags []string
	// TODO
}

func NewSourceFromCLIFlags(flags []string) *SourceCLI {
	return &SourceCLI{flags: flags}
}

func (sc *SourceCLI) FromName(name string) MultiSingleValueSource {
	return nil // TODO
}
func (sc *SourceCLI) FromNameAndShorthand(name, shorthand string) MultiSingleValueSource {
	return nil // TODO
}
func (sc *SourceCLI) FromPositionalArg(position int) MultiSingleValueSource {
	return nil // TODO
}
