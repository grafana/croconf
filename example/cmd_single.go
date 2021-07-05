package main

import (
	"go.k6.io/croconf"
	"go.k6.io/croconf/example/config"
)

func getSingleValCommand(
	configManager *croconf.Manager, globalConf *config.GlobalConfig,
	cliSource *croconf.SourceCLI, envVarsSource *croconf.SourceEnvVars,
) SubCommand {
	var singleTestValue string

	return SubCommand{
		Command: "single",
		AddConfigOptions: func() error {
			configManager.AddField(
				croconf.NewStringField(
					&singleTestValue,
					croconf.DefaultStringValue("foobar"),
					envVarsSource.From("SIMPLE_TEST_VAL_DEPRECATED"), // TODO: warn for deprecated
					envVarsSource.From("SIMPLE_TEST_VAL"),
					cliSource.FromNameAndShorthand("test", "t"),
				),
				croconf.WithDescription("just a simple test value outside of a struct, but still not global"),
			)

			return configManager.Consolidate()
		},
		Run: func() error {
			dumpField(configManager, &singleTestValue)
			return nil
		},
	}
}
