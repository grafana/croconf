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
)

func main() {
	cliSource := croconf.NewSourceFromCLIFlags(os.Args)
	envVarsSource := croconf.NewSourceFromEnv(os.Environ())

	globalConf, err := NewGlobalConfig(cliSource, envVarsSource)
	if err != nil {
		log.Fatal(err)
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
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		log.Fatal(err)
	}
	jsonSource, err := croconf.NewJSONSource(jsonConfigContents)
	if err != nil {
		log.Fatal(err)
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

	fmt.Println()

	vusMeta := scriptConf.cm.Field(&scriptConf.VUs)
	if vusMeta.HasBeenSet() {
		fmt.Printf("Field VUs was manually set by source '%s'\n", vusMeta.SourceOfValue().GetName())
	} else {
		fmt.Printf("Field VUs was using the default value\n")
	}
}
