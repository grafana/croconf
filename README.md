<h1 align="center">croconf</h1>
<h4 align="center">A flexible and composable configuration library for Go</h4>

## Why?

We know that there are plenty of [other Go configuration](https://github.com/avelino/awesome-go#configuration) and [CLI libraries](https://github.com/avelino/awesome-go#standard-cli) out there already - _insert [obligatory xkcd](https://xkcd.com/927/)_... :sweat_smile:  Unfortunately, most (all?) of them suffer from **at least** one of these serious issues and limitations:
1. Difficult to test:
    - e.g. they rely directly on `os.Args()` or `os.Environ()` or some other shared global state
    - can't check what results various inputs will produce without a lot of effort for managing that state
2. Difficult or impossible to extend - some variation of:
    - limited value sources, e.g. they might support CLI flags and env vars, but not JSON or YAML
    - you can't easily write your own custom first-class option types or value sources
    - the value sources are not layered, values from different sources may be difficult or impossible to merge automatically
3. Untyped and reflection-heavy:
    - they fail at run-time instead of compile-time
    - e.g. your app panics because the type for some infrequently used and not very well tested option doesn't implement `encoding.TextUnmarshaler`
    - struct tags are used for _everything_ :scream:
    - alternatively, you may have to do a ton of type assertions deep in your codebase
4. Un-queriable:
    - there is no metadata about the final consolidated config values
    - you cannot know if a certain option was set by the user or if its default value was used
    - you may have to rely on `null`-able or other custom wrapper types for such information
5. Too `string`-y:
    - you have to specify string IDs (e.g. CLI flag names, environment variable names, etc.) multiple times
    - a typo in only some of these these strings might go unnoticed for a long while or cause a panic

The impetus for croconf was [k6](https://github.com/k6io/k6)'s very complicated configuration. We have a lot of options and most options have _at least_ 5 hierarchical value sources: their default values, JSON config, exported `options` in the JS scripts, environment variables, and CLI flag values. Some options have more... :sob:

We currently use several Go config libraries and a lot of glue code to manage this, and it's still a frequent source of bugs and heavy technical debt. As far as we know, no single other existing Go configuration library is sufficient to cover all of our use cases well. And, from what we can see, these issues are only partially explained by Go's weak type system...

So when we tried to find a Go config library that avoids all of these problems and couldn't, croconf was born! :tada:

## Architecture

> ### ⚠️ croconf is still in the "proof of concept" stage
>
> The library is not yet ready for production use. It has bugs, not all features are finished, comments and tests are spotty, and the module structure and type names are expected to change a lot in the coming weeks.

In short, croconf shouldn't suffer from any of the issues :arrow_up:, hopefully without introducing any new ones! :fingers_crossed: It should be suitable for any size of a Go project - from the simplest toy project, to the most complicated CLI application and everything in-between!

Some details about croconf's API design
- it uses type safe, uses plain old Go values for the config values
- works for standalone values as well as `struct` properties
- everything about a config field is defined in a single place, no `string` identifier has to ever be written more than once
- after consolidating the config values, you can query which config source was responsible for setting a specific value (or if the default value was set)
- batteries included, while at the same time completely extensible:
    - built-in frontends for all native Go types, incl. `encoding.TextUnmarshaler` and slices
    - support for CLI flags, environment variables and JSON options (and others in the future) out of the box, with zero dependencies
    - none of the built-in types are special, you can easily add custom value types and config sources by implementing a few of the small well-defined interfaces in [`types.go`](https://github.com/k6io/croconf/blob/main/types.go)
- no `unsafe` and no magic :sparkles:
- no `reflect` and no type assertions needed for user-facing code (both are used very sparingly internally in the library)

These nice features and guarantees are achieved because of the type-safe lazy bindings between value destinations and source paths that croconf uses. The configuration definition just defines the source bindings for every value, the actual resolving is done as a subsequent step.

## Example

```go
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
		croconf.WithDescription("number of virtual users"),
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

	jsonResult, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		log.Fatalf("error marshaling JSON: %s", err)
	}
	fmt.Fprint(os.Stdout, string(jsonResult))
}
```

This was a relatively simple example taken from [here](https://github.com/k6io/croconf/blob/main/examples/croconf-simple-struct-example/main.go), and it still manages to combine 4 config value sources! For other examples, take a look in the [`examples` folder](https://github.com/k6io/croconf/tree/main/examples) in this repo.

## Origins of name

croconf comes from _croco_-dile _conf_-iguration. So, :crocodile: not :croatia: :smile: And in the tradition set by [k6](https://github.com/k6io/k6), if we don't like it, we might decide to abbreviate it to `c6` later... :sweat_smile:

## Remaining tasks

As mentioned above, this library is still in the proof-of-concept stage. It is usable for toy projects and experiments, but it is very far from production-ready. These are some of the remaining tasks:
- Refactor module structure and type names
- More value sources (e.g. TOML, YAML, INI, etc.) and improvements in the current ones
- Add built-in support for all Go basic and common stdlib types and interfaces
- Code comments and linter fixes
- Fix bugs and write **a lot** more tests
- Documentation and examples
- Better (more user-friendly) error messages
- An equivalent to [cobra](https://github.com/spf13/cobra) or [kong](https://github.com/alecthomas/kong), a wrapper for CLI application frameworks that is able to handle CLI sub-commands, shell autocompletion, etc.
    - _currently only toy PoC for this concept exists in [`examples/croconf-complex-example/`](https://github.com/k6io/croconf/tree/main/examples/croconf-complex-example)_
- Add drop-in support for marshaling config structs (e.g. to JSON) with the same format they were unmarshaled from.
- Be able to emit errors on unknown CLI flags, JSON options, etc.
