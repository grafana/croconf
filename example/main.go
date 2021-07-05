package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"

	"go.k6.io/croconf"
	"go.k6.io/croconf/example/config"
)

func main() {
	configManager := croconf.NewManager()
	cliSource := croconf.NewSourceFromCLIFlags(os.Args[1:])
	envVarsSource := croconf.NewSourceFromEnv(os.Environ())

	globalConf := config.NewGlobalConfig(configManager, cliSource, envVarsSource)

	if err := configManager.Consolidate(); err != nil {
		log.Fatal(err)
	}

	// At this point, there are plenty of unknown options still, but we should
	// at least know which sub-command we need to execute, and we should be able
	// to handle things like --help

	// TODO: obviously something better
	if globalConf.SubCommand == "run" {
		runCommand(configManager, globalConf, cliSource, envVarsSource)
	} else if globalConf.ShowHelp {
		fmt.Println(configManager.GetHelpText()) //nolint:forbidigo
	} else {
		log.Fatalf("unknown sub-command %s, see options with --help", globalConf.SubCommand)
	}
}

//nolint:forbidigo
func runCommand(
	configManager *croconf.Manager, globalConf *config.GlobalConfig,
	cliSource *croconf.SourceCLI, envVarsSource *croconf.SourceEnvVars,
) {
	jsonConfigContents, err := ioutil.ReadFile(globalConf.JSONConfigPath)
	if err != nil {
		if configManager.Field(&globalConf.JSONConfigPath).HasBeenSetFromSource() {
			// If this was explicitly set, treat any failure to open it as a fatal error
			log.Fatal(err)
		}
		if !errors.Is(err, fs.ErrNotExist) {
			// if we're using the default log config location, warn on any errors except "file not found"
			log.Printf("warning! could not open config.json file: %s", err)
		}
	}
	jsonSource := croconf.NewJSONSource(jsonConfigContents)
	scriptConf := config.NewScriptConfig(configManager, globalConf, cliSource, envVarsSource, jsonSource)
	if err := configManager.Consolidate(); err != nil {
		log.Fatal(err)
	}

	if scriptConf.ShowHelp {
		fmt.Println(configManager.GetHelpText()) //nolint:forbidigo
		return
	}

	// TODO error out if we see unknown CLI flags or JSON options

	// And finally, we should be able to marshal and dump the consolidated config
	jsonResult, err := json.MarshalIndent(scriptConf, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonResult))

	fmt.Println()

	dumpField(configManager, &scriptConf.JSONConfigPath)
	dumpField(configManager, &scriptConf.VUs)
	dumpField(configManager, &scriptConf.Scenarios1)
	dumpField(configManager, &scriptConf.Scenarios2)
	dumpField(configManager, &scriptConf.DNS.TTL)
	dumpField(configManager, &scriptConf.DNS.Server)
	dumpField(configManager, &scriptConf.TinyArr)
	dumpField(configManager, &scriptConf.Throw)
}

//nolint:forbidigo
func dumpField(cm *croconf.Manager, field interface{}) {
	// TODO: get the name from the field?
	jsonResult, err := json.Marshal(field)
	if err != nil {
		log.Fatal(err)
	}

	fieldMeta := cm.Field(field)
	if fieldMeta.HasBeenSetFromSource() {
		binding := fieldMeta.LastBindingFromSource()
		fmt.Printf(
			"Field %s was manually set by source '%s' (field %s) with value '%s'\n",
			fieldMeta.Name, binding.Source().GetName(), binding.BoundName(), jsonResult,
		)
	} else {
		fmt.Printf(
			"Field %s was using the default/non-source value of '%s'\n",
			fieldMeta.Name, jsonResult,
		)
	}
}
