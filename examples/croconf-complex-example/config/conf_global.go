package config

import (
	"go.k6.io/croconf"
)

type GlobalConfig struct {
	Verbose        bool
	JSONConfigPath string
	// TODO: other global or runtime options...
}

func NewGlobalConfig(
	cm *croconf.Manager, cliSource *croconf.SourceCLI, envVarsSource *croconf.SourceEnvVars,
) *GlobalConfig {
	conf := &GlobalConfig{}

	cm.AddField(
		croconf.NewStringField(
			&conf.JSONConfigPath,
			croconf.DefaultStringValue("./config.json"),
			envVarsSource.From("K6_CONFIG"),
			cliSource.FromNameAndShorthand("config", "c"),
		),
		croconf.WithDescription("path to k6 JSON config file"),
	)

	cm.AddField(
		croconf.NewBoolField(
			&conf.Verbose,
			envVarsSource.From("K6_VERBOSE"),
			cliSource.FromNameAndShorthand("verbose", "v"),
		),
		croconf.WithDescription("enable verbose logging"),
	)

	// TODO: add the other global options and properties

	return conf
}
