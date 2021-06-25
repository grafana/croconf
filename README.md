# croconf

A flexible and composable configuration library for Go that doesn't suck

### Ned's spec for Go configuration which doesn't suck:

1. Fully testable: there is no relying on globals or os directly, everything is passed as parameters (e.g. it receives os.Environ() and os.Args, it doesn't directly access them).
1. Supports layered configs: users can construct hierarchies of json/yaml/toml/env vars/CLI flags/etc., and the library will merge them
1. Uses normal and simple Go types:
    - the end value should be a plain old Go struct with plain old Go types (e.g. string, int, etc.)
    - at the same time, users should have a way to query and access metadata to answer questions like "Has field X been changed?", "What is the default value of field Y?", etc.
    - there will be no custom types to check if an entry was set (i.e. no null.Int.Valid BS...)
    - the final consolidation result is a plain Go struct and a separate metadata layer allows users to reason about the config and answer the questions above
    - custom types will only be needed for complex options (and the library will have nicely defined interfaces for supporting custom types)
1. This needs to be composable, in all three dimensions:
    - config values can be consolidated between multiple config layers (e.g. CLI flag overwrites env. var which overwrites JSON option, etc.)
    - configs can be combined (e.g. if I have configs for type A struct { ... } and type B struct { ...}, this should also be easy to make into a valid config: type C struct {*A; *B}
    - a config can contain another config as a property, i.e. you should be able to nest configs, and one config can be encapsulated in a single property of the other
1. Everything is as type safe and compile-time-error-able as possible:
    - static go types and interfaces >> type assertions >> reflection
    - we won't use struct tags! type-safe methods/properties >>> struct tags
1. Batteries built-in (e.g. support for JSON, env vars, CLI flags, basic data types), but completely extensible
    - Supports validation, has to have user-friendly error messages (without Go implementation details)
    - Supports warnings for things like deprecated variables
1. An easy way to marshal the whole consolidated config, e.g. to a JSON file. Ideally, we should be able to specify whether we want only the changed values, or all of the values (incl. any default ones).
1. Stretch goal: the metadata should be rich enough so that a whole application framework like cobra can be built on top of it, including generation of man pages and auto-completion


### Proposed TODO:
0. Figure out a usable Go API (e.g. with initial support for just a few Go types like `string`, `int64` and `bool` that satisfies the criteria :arrow_up: :sweat_smile:
1. Write a PoC with some tests and mock real-life usage examples
2. Iterate and expand on :arrow_up:
3. Support all types (incl. custom types) and multiple sources
4. Polish, set up GitHub Actions CI, etc.
5. Profit

