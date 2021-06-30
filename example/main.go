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
	cliSource := croconf.NewSourceFromCLIFlags(os.Args[1:])
	envVarsSource := croconf.NewSourceFromEnv(os.Environ())

	globalConf, err := NewGlobalConfig(cliSource, envVarsSource)
	if err != nil {
		log.Fatal(err)
	}

	// At this point, there are plenty of unknown options still, but we should
	// at least know which sub-command we need to execute, and we should be able
	// to handle things like --help

	// TODO: obviously something better
	if globalConf.SubCommand == "run" {
		runCommand(cliSource, envVarsSource, globalConf)
	} else {
		log.Fatalf("unknown sub-command %s", globalConf.SubCommand)
	}
}

//nolint:forbidigo
func runCommand(
	cliSource *croconf.SourceCLI,
	envVarsSource *croconf.SourceEnvVars,
	globalConf *GlobalConfig,
) {
	jsonConfigContents, err := ioutil.ReadFile(globalConf.JSONConfigPath)
	if err != nil {
		if globalConf.cm.Field(&globalConf.JSONConfigPath).ValueSource() != nil {
			// If this was explicitly set, treat any failure to open it as a fatal error
			log.Fatal(err)
		}
		if !errors.Is(err, fs.ErrNotExist) {
			// if we're using the default log config location, warn on any errors except "file not found"
			log.Printf("warning! could not open config.json file: %s", err)
		}
	}
	jsonSource := croconf.NewJSONSource(jsonConfigContents)
	scriptConf, err := NewScriptConfig(globalConf, cliSource, envVarsSource, jsonSource)
	if err != nil {
		log.Fatal(err)
	}

	// And finally, we should be able to marshal and dump the consolidated config
	jsonResult, err := json.MarshalIndent(scriptConf, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonResult))

	fmt.Println()

	dumpField(scriptConf.cm, &scriptConf.VUs, "VUs")
	dumpField(scriptConf.cm, &scriptConf.Scenarios1, "Scenarios1")
	dumpField(scriptConf.cm, &scriptConf.Scenarios2, "Scenarios2")
	dumpField(scriptConf.cm, &scriptConf.DNS.TTL, "DNS.TTL")
	dumpField(scriptConf.cm, &scriptConf.DNS.Server, "DNS.Server")
}

//nolint:forbidigo
func dumpField(cm *croconf.Manager, field interface{}, fieldName string) {
	// TODO: get the name from the field?
	jsonResult, err := json.Marshal(field)
	if err != nil {
		log.Fatal(err)
	}

	fieldMeta := cm.Field(field)
	if source := fieldMeta.ValueSource(); source != nil {
		fmt.Printf(
			"Field %s was manually set by source '%s' with value '%s'\n",
			fieldName, source.GetName(), jsonResult,
		)
	} else {
		fmt.Printf(
			"Field %s was using the default value of '%s'\n",
			fieldName, jsonResult,
		)
	}
}
