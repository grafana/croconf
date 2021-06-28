package main

import (
	"go.k6.io/croconf"
)

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

	cm.AddField(
		croconf.NewStringField(
			&conf.SubCommand,
			croconf.DefaultStringValue("run"),
			//cliSource.FromPositionalArg(1),
		),
		/*
			croconf.IsRequired(),
			croconf.WithDescription("k6 sub-command"),
			croconf.WithValidator(func() error {
				// TODO: validate
				if conf.SubCommand != "run" || conf.SubCommand != "cloud" {
					return errors.New("foo")
				}
			}),
		*/
	)

	cm.AddField(
		croconf.NewStringField(
			&conf.JSONConfigPath,
			croconf.DefaultStringValue("~/.config/loadimpact/k6/config.json"),
			envVarsSource.From("K6_CONFIG"),
			cliSource.FromNameAndShorthand("--config", "-c"),
		),
		/*
			croconf.WithDescription("path to k6 JSON config file"),
			croconf.WithValidator(func() error {
				// TODO: validate
				if conf.SubCommand != "run" || conf.SubCommand != "cloud" {
					return errors.New("foo")
				}
			}),
		*/
	)

	// TODO: add the other options and properties

	if err := cm.Consolidate(); err != nil {
		return nil, err
	}

	return conf, nil
}
