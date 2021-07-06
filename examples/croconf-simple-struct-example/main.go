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

	"go.k6.io/croconf"
)

// SimpleConfig is a normal Go struct with plain Go property types.
type SimpleConfig struct {
	RPPs int64
	DNS  struct {
		Server net.IP // type that implements encoding.TextUnmarshaler
		// ... more nested fields
	}
	// ... more config fields...
}

// NewScriptConfig defines the sources and metadata for every config field.
func NewScriptConfig(
	cm *croconf.Manager, cliSource *croconf.SourceCLI,
	envVarsSource *croconf.SourceEnvVars, jsonSource *croconf.SourceJSON,
) *SimpleConfig {
	conf := &SimpleConfig{}

	cm.AddField(
		croconf.NewInt64Field(
			&conf.RPPs,
			jsonSource.From("rps"),
			envVarsSource.From("APP_RPS"),
			cliSource.FromNameAndShorthand("rps", "r"),
			// ... more bindings - every field can have as many or as few as needed
		),
		croconf.WithDescription("requests per second"),
		croconf.IsRequired(),
		// ... more field options like validators, meta-information, etc.
	)

	cm.AddField(
		croconf.NewTextBasedField(
			&conf.DNS.Server,
			croconf.DefaultStringValue("8.8.8.8"),
			jsonSource.From("dns").From("server"),
			envVarsSource.From("APP_DNS_SERVER"),
		),
		croconf.WithDescription("server for DNS queries"),
	)

	// ... more fields

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
