package main

import (
	"go.k6.io/croconf"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/executor"
)

type ScriptConfig struct {
	cm *croconf.Manager
	*GlobalConfig

	UserAgent string
	VUs       int64

	Duration Duration

	Scenarios1 lib.ScenarioConfigs

	Scenarios2 lib.ScenarioConfigs

	// TODO: have a sub-config
}

func NewScriptConfig(
	globalConf *GlobalConfig,
	cliSource *croconf.SourceCLI,
	envVarsSource *croconf.SourceEnvVars,
	jsonSource *croconf.SourceJSON,
) (*ScriptConfig, error) {
	cm := croconf.NewManager()
	conf := &ScriptConfig{GlobalConfig: globalConf, cm: cm} // TODO: somehow save the sources in the struct as well?

	cm.AddField(
		croconf.NewStringField(
			&conf.UserAgent,
			croconf.DefaultStringValue("croconf example demo v0.0.1 (https://k6.io/)"),
			jsonSource.From("userAgent"),
			envVarsSource.From("K6_USER_AGENT"),
			cliSource.FromName("user-agent"),
			// TODO: figure this out...
			// croconf.WithDescription("user agent for http requests")
		),
	)

	cm.AddField(croconf.NewInt64Field(
		&conf.VUs,
		croconf.DefaultInt64Value(1),
		jsonSource.From("vus"),
		envVarsSource.From("K6_VUS"),
		cliSource.FromNameAndShorthand("vus", "u"),
		// croconf.WithDescription("number of virtual users") // TODO
	))

	cm.AddField(croconf.NewTextBasedField(
		&conf.Duration,
		croconf.DefaultStringValue("1s"),
		jsonSource.From("duration"),
		envVarsSource.From("K6_DURATION"),
		cliSource.FromNameAndShorthand("duration", "d"),
	))

	// This is one way to add a custom field in a type-safe manner:
	cm.AddField(croconf.NewCustomField(
		&conf.Scenarios1,
		croconf.DefaultCustomValue(func() {
			conf.Scenarios1 = lib.ScenarioConfigs{
				lib.DefaultScenarioName: executor.NewPerVUIterationsConfig(lib.DefaultScenarioName),
			}
		}),
		jsonSource.From("scenarios1").To(&conf.Scenarios1),
	))
	// This is another way to do it, with a small custom type:
	cm.AddField(&scenariosField{dest: &conf.Scenarios2, source: jsonSource, id: "scenarios2"})

	// TODO: add the other options and actually process and consolidate the
	// config values and handle any errors... Here we probably want to error out
	// if we see unknown CLI flags or JSON options

	// TODO: automatically do this on Consolidate()?
	if err := cliSource.Parse(); err != nil {
		return nil, err
	}

	if err := cm.Consolidate(); err != nil {
		return nil, err
	}

	return conf, nil
}

type scenariosField struct {
	dest   *lib.ScenarioConfigs
	source *croconf.SourceJSON
	id     string
	wasSet bool
}

func (sf *scenariosField) Destination() interface{} {
	return sf.dest
}

func (sf *scenariosField) Consolidate() []error {
	*sf.dest = lib.ScenarioConfigs{
		lib.DefaultScenarioName: executor.NewPerVUIterationsConfig(lib.DefaultScenarioName),
	}

	jsonVal, ok := sf.source.Lookup(sf.id)
	if !ok {
		return nil
	}

	if err := sf.dest.UnmarshalJSON(jsonVal); err != nil {
		return []error{err}
	}

	return nil
}
func (sf *scenariosField) ValueSource() croconf.Source {
	if sf.wasSet {
		return sf.source
	}
	return nil
}
