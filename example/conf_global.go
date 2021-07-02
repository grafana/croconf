package main

import (
	"fmt"

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
	cm.AddSource(envVarsSource)
	cm.AddSource(cliSource)
	conf := &GlobalConfig{cm: cm}

	cm.AddField(
		croconf.NewStringField(
			&conf.SubCommand,
			cliSource.FromPositionalArg(1),
		),
		croconf.WithDescription("k6 sub-command"),
		croconf.IsRequired(),
		croconf.WithValidator(func() error {
			if conf.SubCommand != "run" && conf.SubCommand != "cloud" {
				return fmt.Errorf("invalid sub-command %s", conf.SubCommand)
			}
			return nil
		}),
	)

	cm.AddField(
		croconf.NewStringField(
			&conf.JSONConfigPath,
			croconf.DefaultStringValue("./config.json"),
			envVarsSource.From("K6_CONFIG"),
			cliSource.FromNameAndShorthand("config", "c"),
		),
		croconf.WithDescription("path to k6 JSON config file"),
	)

	// TODO: add the other options and properties

	if err := cm.Consolidate(); err != nil {
		return nil, err
	}

	return conf, nil
}
