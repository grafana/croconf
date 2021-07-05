package config

import (
	"fmt"

	"go.k6.io/croconf"
)

type GlobalConfig struct {
	Verbose        bool
	JSONConfigPath string
	// TODO: other global or runtime options...

	// TODO: move out of here
	SubCommand string // run, cloud, inspect, archive, etc.
	ShowHelp   bool
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

	// TODO: move to CLI framework helper
	cm.AddField(
		croconf.NewStringField(
			&conf.SubCommand,
			cliSource.FromPositionalArg(1),
		),
		croconf.WithDescription("k6 sub-command"),
		croconf.WithValidator(func() error {
			if conf.SubCommand != "" && conf.SubCommand != "run" && conf.SubCommand != "cloud" {
				return fmt.Errorf("invalid sub-command %s", conf.SubCommand)
			}
			return nil
		}),
	)

	cm.AddField(
		croconf.NewBoolField(
			&conf.ShowHelp,
			cliSource.FromNameAndShorthand("help", "h"),
		),
		croconf.WithDescription("show help information"),
	)

	return conf
}
