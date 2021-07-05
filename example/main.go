package main

import (
	"log"
	"os"

	"go.k6.io/croconf"
	"go.k6.io/croconf/example/config"
)

func main() {
	cliSource := croconf.NewSourceFromCLIFlags(os.Args[1:])
	envVarsSource := croconf.NewSourceFromEnv(os.Environ())
	configManager := croconf.NewManager(
		croconf.WithDefaultSourceOfFieldNames(envVarsSource),
	)

	globalConf := config.NewGlobalConfig(configManager, cliSource, envVarsSource)

	subCommands := getSubcommands(configManager, globalConf, cliSource, envVarsSource)

	handler, err := GetSubcommandHandler(configManager, cliSource, subCommands, []croconf.StringValueBinder{
		envVarsSource.From("K6_SUB_COMMAND"),
		cliSource.FromPositionalArg(1),
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := configManager.Consolidate(); err != nil {
		log.Fatal(err)
	}

	if err := handler(); err != nil {
		log.Fatal(err)
	}
}

func getSubcommands(
	configManager *croconf.Manager, globalConf *config.GlobalConfig,
	cliSource *croconf.SourceCLI, envVarsSource *croconf.SourceEnvVars,
) []SubCommand {
	return []SubCommand{
		runCommand(configManager, globalConf, cliSource, envVarsSource),
		getSingleValCommand(configManager, globalConf, cliSource, envVarsSource),
	}
}
