package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"go.k6.io/croconf"
)

func main() {
	cliSource := croconf.NewSourceFromCLIFlags(os.Args)
	envVarsSource := croconf.NewSourceFromEnv(os.Environ())

	globalConf, err := NewGlobalConfig(cliSource, envVarsSource)
	if err != nil {
	}

	// At this point, there are plenty of unknown options still, but we should
	// at least know which sub-command we need to execute, and we should be able
	// to handle things like --help

	// TODO: obviosuly something better
	if globalConf.SubCommand == "run" || true /* TODO: remove after we actually populate the option */ {
		runCommand(cliSource, envVarsSource, globalConf)
	}
}

func runCommand(
	cliSource *croconf.SourceCLI,
	envVarsSource *croconf.SourceEnvVars,
	globalConf *GlobalConfig,
) {
	jsonConfigContents, err := ioutil.ReadFile(globalConf.JSONConfigPath)
	if err != nil {
		// TODO: handle error
	}
	jsonSource, err := croconf.NewJSONSource(jsonConfigContents)
	if err != nil {
		// TODO: handle error
	}

	scriptConf, err := NewScriptConfig(globalConf, cliSource, envVarsSource, jsonSource)
	if err != nil {
		// TODO: handle error
	}

	// And finally, we should be able to marshal and dump the consolidated config
	jsonResult, err := json.Marshal(scriptConf)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonResult))
}
