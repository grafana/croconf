package croconf

type SourceEnvVars struct {
	environ []string
	// TODO
}

func NewSourceFromEnv(environ []string) *SourceEnvVars {
	return &SourceEnvVars{environ: environ}
}

func (sev *SourceEnvVars) From(name string) MultiSingleValueSource {
	return nil // TODO
}
