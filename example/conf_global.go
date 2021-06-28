package main

import "go.k6.io/croconf"

type GlobalConfig struct {
	cm *croconf.Manager

	// TODO: embed the CLI and env var sources?

	SubCommand     string // run, cloud, inspect, archive, etc.
	JSONConfigPath string
	// TODO: other global or runtime options...
}

func NewGlobalConfig(
	cliSource *croconf.SourceCLI,
	envVarsSource *croconf.SourceEnvVars,
) (*GlobalConfig, error) {
	cm := croconf.NewManager()
	conf := &GlobalConfig{cm: cm}

	cm.StringField(
		&conf.SubCommand,
		croconf.DefaultStringValue("run"),
		cliSource.FromPositionalArg(1),
	)

	cm.StringField(
		&conf.JSONConfigPath,
		croconf.DefaultStringValue("~/.config/loadimpact/k6/config.json"),
		envVarsSource.From("K6_CONFIG"),
		cliSource.FromNameAndShorthand("--config", "-c"),
	)

	// TODO: add the other options and actually process and consolidate the
	// config values and handle any errors

	return conf, nil
}
