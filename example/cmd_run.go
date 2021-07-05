package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"

	"go.k6.io/croconf"
	"go.k6.io/croconf/example/config"
)

//nolint:forbidigo
func runCommand(
	configManager *croconf.Manager, globalConf *config.GlobalConfig,
	cliSource *croconf.SourceCLI, envVarsSource *croconf.SourceEnvVars,
) SubCommand {
	var jsonSource *croconf.SourceJSON
	var scriptConf *config.ScriptConfig

	return SubCommand{
		Command: "run",
		AddConfigOptions: func() error {
			// TODO: actually do this after we've applied the options, the
			// bindings don't depend on us having actually read the JSON config
			// file contents, so no need to do it now
			jsonConfigContents, err := ioutil.ReadFile(globalConf.JSONConfigPath)

			// If this was explicitly set, treat any failure to open it as a
			// fatal error. If we're using the default log config location, do
			// not consider "file not found" an error.
			if err != nil && (configManager.Field(&globalConf.JSONConfigPath).HasBeenSetFromSource() ||
				!errors.Is(err, fs.ErrNotExist)) {
				return fmt.Errorf("could not open json config file: %w", err)
			}

			jsonSource = croconf.NewJSONSource(jsonConfigContents)
			scriptConf = config.NewScriptConfig(configManager, globalConf, cliSource, envVarsSource, jsonSource)

			// TODO: error out if we see unknown CLI flags or JSON options

			return configManager.Consolidate()
		},
		Run: func() error {
			// And finally, we should be able to marshal and dump the consolidated config
			jsonResult, err := json.MarshalIndent(scriptConf, "", "    ")
			if err != nil {
				return err
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

			return nil
		},
	}
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
