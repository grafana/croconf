package config

import (
	"go.k6.io/croconf"
	"go.k6.io/croconf/example/types"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/executor"
)

type ScriptConfig struct {
	*GlobalConfig
	// TODO: json.Marshaler

	UserAgent string
	VUs       int64

	Throw bool

	Duration types.Duration

	DNS struct {
		TTL    types.Duration
		Server string
	}

	Tiny    int8
	TinyArr []int8

	Scenarios1 lib.ScenarioConfigs

	Scenarios2 lib.ScenarioConfigs
}

// TODO: split apart in multiple functions, one per section?
func NewScriptConfig( //nolint: funlen
	cm *croconf.Manager, globalConf *GlobalConfig,
	cliSource *croconf.SourceCLI, envVarsSource *croconf.SourceEnvVars, jsonSource *croconf.SourceJSON,
) *ScriptConfig {
	conf := &ScriptConfig{
		GlobalConfig: globalConf,
		// TODO: implement something like this:
		// Marshaler:    jsonSource.NewMarshaller(),
	}

	cm.AddField(
		croconf.NewStringField(
			&conf.UserAgent,
			croconf.DefaultStringValue("croconf example demo v0.0.1 (https://k6.io/)"),
			jsonSource.From("userAgent"),
			envVarsSource.From("K6_USER_AGENT"),
			cliSource.FromName("user-agent"),
		),
		croconf.WithDescription("user agent for http requests"),
	)

	cm.AddField(
		croconf.NewInt64Field(
			&conf.VUs,
			croconf.DefaultIntValue(1),
			jsonSource.From("vus"),
			envVarsSource.From("K6_VUS"),
			cliSource.FromNameAndShorthand("vus", "u"),
		),
		croconf.WithDescription("number of virtual users"),
		croconf.IsRequired(),
	)

	cm.AddField(
		croconf.NewBoolField(
			&conf.Throw,
			// TODO: croconf.DefaultBoolValue(false),
			jsonSource.From("throw"),
			envVarsSource.From("K6_THROW"),
			cliSource.FromNameAndShorthand("throw", "w"),
		),
		croconf.WithDescription("throw warnings (like failed http requests) as errors"),
	)

	cm.AddField(
		croconf.NewTextBasedField(
			&conf.Duration,
			croconf.DefaultStringValue("1s"),
			jsonSource.From("duration"),
			envVarsSource.From("K6_DURATION"),
			cliSource.FromNameAndShorthand("duration", "d"),
		),
		croconf.WithDescription("test duration"),
	)

	// Properties of a nested struct (without a pointer!)
	cm.AddField(croconf.NewTextBasedField(
		&conf.DNS.TTL,
		croconf.DefaultStringValue("10m"),
		jsonSource.From("dns").From("ttl"),
	))

	cm.AddField(croconf.NewStringField(
		&conf.DNS.Server,
		croconf.DefaultStringValue("8.8.8.8"),
		jsonSource.From("dns").From("server"),
	))

	cm.AddField(croconf.NewInt8Field(
		&conf.Tiny,
		croconf.DefaultIntValue(1),
		jsonSource.From("tiny"),
		envVarsSource.From("K6_TINY"),
		cliSource.FromName("tiny"),
	))

	cm.AddField(croconf.NewInt8SliceField(
		&conf.TinyArr,
		jsonSource.From("tinyArr"),
		envVarsSource.From("K6_TINY_ARR"),
		cliSource.FromNameAndShorthand("tiny-arr", "a"),
		// TODO: other sources and defaults
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

	// TODO: add the other options

	return conf
}
