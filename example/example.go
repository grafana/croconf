package main

import (
	"encoding/json"
	"fmt"

	"go.k6.io/croconf"
)

type Config struct {
	cm *croconf.Manager

	UserAgent string
	VUs       int64

	// TODO: have a sub-config
}

func NewConfig(
	jsonSource *croconf.SourceJSON,
	envVarsSource *croconf.SourceEnvVars,
	cliSource *croconf.SourceCLI,
) *Config {
	cm := croconf.NewManager()
	conf := &Config{cm: cm} // TODO: somehow save the sources in the struct as well?

	cm.StringField(
		&conf.UserAgent, "croconf example demo v0.0.1 (https://k6.io/)",
		jsonSource.From("userAgent"),
		envVarsSource.From("K6_USER_AGENT"),
		cliSource.FromName("--user-agent"),
		// TODO: figure this out...
		// croconf.WithDescription("user agent for http requests")
	)

	cm.Int64Field(
		&conf.VUs, 1,
		jsonSource.From("vus"),
		envVarsSource.From("K6_VUS"),
		cliSource.FromNameAndShorthand("--vus", "-u"),
		// croconf.WithDescription("number of virtual users") // TODO
	)

	return conf
}

func main() {
	jsonSource := &croconf.SourceJSON{}
	envVarsSource := &croconf.SourceEnvVars{}
	cliSource := &croconf.SourceCLI{}

	conf := NewConfig(jsonSource, envVarsSource, cliSource)
	// at this point, we haven't loaded any values, the various sources are
	// empty, but we should have all of the metadata for go run . --help to work
	// (i.e. to implement the cobra-like library on top of our config, with help
	// and subcommands and everything)

	// TODO: after we have the metadata, we actually have to load the config data,
	// somewhat like this (while actually handling errors):

	// envVarsSource.Parse(os.Environ())
	// cliSource.Parse(os.Args)

	// TODO: figure out the config JSON path from its default value in conf and
	// potentially influenced by the an environment var or a CLI flags

	// jsonConfigContents, err := ioutil.ReadFile(consolidatedPathToJSONConfig)
	// jsonSource.Parse(jsonConfigContents)

	// err := conf.ConsolidateValuesFromSources()
	// TODO: handle err
	// err := conf.Validate()
	// TODO: handle err

	// And finally, we should be able to marshal and dump the consolidated config
	jsonResult, err := json.Marshal(conf)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonResult))
}
