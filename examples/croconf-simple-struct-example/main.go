package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"

	"go.k6.io/croconf"
)

type MyCustomType struct {
	foo, bar string
}

func MyCustomTypeFromStr(s string) (MyCustomType, error) {
	parts := strings.SplitN(s, ":", 2)
	res := MyCustomType{foo: parts[0]}
	if len(parts) > 1 {
		res.bar = parts[1]
	}
	return res, nil
}

// SimpleConfig is a normal Go struct with plain Go property types.
type SimpleConfig struct {
	RPS int64
	DNS struct {
		Server net.IP // type that implements encoding.TextUnmarshaler
		// ... more nested fields
	}
	MyCustom MyCustomType
	// ... more config fields...
}

// NewScriptConfig defines the sources and metadata for every config field.
func NewScriptConfig(
	cm *croconf.Manager, cliSource *croconf.SourceCLI,
	envVarsSource *croconf.SourceEnvVars, jsonSource *croconf.SourceJSON,
) *SimpleConfig {
	conf := &SimpleConfig{}

	cm.AddField(
		croconf.NewField(&conf.RPS).
			WithDefault(100).
			WithBinding(jsonSource.From("rps").ToInt64()).
			WithBinding(envVarsSource.From("APP_RPS").ToInt64()).
			WithBinding(cliSource.FromNameAndShorthand("rps", "r").ToInt64()),
	)

	cm.AddField(
		croconf.NewField(&conf.MyCustom).
			WithBinding(
				croconf.StringToCustomType(envVarsSource.From("MYCUSTOM").ToString(), MyCustomTypeFromStr),
			))

	/*
		cm.AddField(
			croconf.NewTextBasedField(
				&conf.DNS.Server,
				croconf.DefaultStringValue("8.8.8.8"),
				jsonSource.From("dns").From("server"),
				envVarsSource.From("APP_DNS_SERVER"),
			),
			croconf.WithDescription("server for DNS queries"),
		)
	*/

	return conf
}

func main() {
	configManager := croconf.NewManager()
	// Manually create config sources - fully testable, no implicit shared globals!
	cliSource := croconf.NewSourceFromCLIFlags(os.Args[1:])
	envVarsSource := croconf.NewSourceFromEnv(os.Environ())
	jsonSource := croconf.NewJSONSource(getJSONConfigContents())

	config := NewScriptConfig(configManager, cliSource, envVarsSource, jsonSource)

	if err := configManager.Consolidate(); err != nil {
		log.Fatalf("error consolidating the config: %s", err)
	}

	printConfig(config) // TODO: something more useful
}

func printConfig(config *SimpleConfig) {
	jsonResult, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		log.Fatalf("error marshaling JSON: %s", err)
	}
	fmt.Fprint(os.Stdout, string(jsonResult))
}

func getJSONConfigContents() []byte {
	// See the croconf-complex-example for how this path can be configured from
	// the CLI flags or environment variables in a multi-step process.
	jsonConfigContents, err := ioutil.ReadFile("./config.json")
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		log.Fatalf("could not open json config file: %s", err)
	}
	return jsonConfigContents
}
