# croconf

A flexible and composable configuration library for Go that doesn't suck

### Ned's spec for Go configuration which doesn't suck:

1. Fully testable: there is no relying on globals or `os` directly, everything is passed as parameters (e.g. it receives `os.Environ()` and `os.Args`, it doesn't directly access them).
2. Supports layered configs: users can construct hierarchies of json/yaml/toml/env vars/CLI flags/etc., and the library will merge them
3. Uses normal and simple Go types:
    - the end value should be a plain old Go `struct` with plain old Go types (e.g. `string`, `int`, etc.)
    - at the same time, users should have a way to query and access metadata to answer questions like "Has field `X` been changed?", "What is the default value of field `Y`?", etc.
    - there will be no custom types to check if an entry was set (i.e. no `null.Int.Valid` BS...)
    - the final consolidation result is a plain Go struct and a separate metadata layer allows users to reason about the config and answer the questions above
    - custom types will only be needed for complex options (and the library will have nicely defined interfaces for supporting custom types)
4. This needs to be composable, in all three dimensions:
    - config values can be consolidated between multiple config layers (e.g. CLI flag overwrites env. var which overwrites JSON option, etc.)
    - configs can be combined (e.g. if I have configs for `type A struct { ... }` and `type B struct { ...}`, this should also be easy to make into a valid config: `type C struct {A; B}`
    - a config can contain another config as a property, i.e. you should be able to nest configs, and one config can be encapsulated in a single property of the other
5. Everything is as type safe and compile-time-error-able as possible:
    - static go types and interfaces >> type assertions >> reflection
    - we won't use struct tags! type-safe methods/properties >>> struct tags
6. Batteries built-in (e.g. support for JSON, env vars, CLI flags, basic data types), but completely extensible
    - Supports validation, has to have user-friendly error messages (without Go implementation details)
    - Supports warnings for things like deprecated variables
7. An easy way to marshal the whole consolidated config, e.g. to a JSON file. Ideally, we should be able to specify whether we want only the changed values, or all of the values (incl. any default ones).
8. Stretch goal: the metadata should be rich enough so that a whole application framework like cobra can be built on top of it, including generation of man pages and auto-completion


### Misc thoughts:
- The building of the final config can be a multi-step process. For example, you may first need to understand which sub-command is going to be used (e.g. `k6 run`, `k6 cloud`, `k6 resume`, etc.), before you actually know _what_ config options are even possible.
- At _some_ of the config building steps, we need to be able to check _some_ config sources for uknown/unused options. For example:
    - at the first step when we're determining the sub-command, we don't care that there will be unknown CLI flags, we expect that
    - at the next step, when we know the sub-command and all of its needed CLI flags and environment variables, an uknown CLI flag should be an error, but an unknown env var shouldn't be.
    - an unknown JSON option might be an error in some places, but for compatibility reasons, a warning in others...
- If the config objects are pointers, and config properties are values in the config structs but passed by pointers to the croconf functions, you have these pros and cons:
    - pro: mostly have a very type safe API without reflection/type assertion
    - pro: you can use the property pointers as keys in the "Has field `X` been changed?" questions
    - con: some config user will be able to modify the config deep in the codebase
    - pro/con: you can copy the config values by just copying the struct, but if you have nested structs by pointer or a `crocon.Manager` (if we stick with that), it will be a big problem...
- Error reporting is tricky... we want it to be as user-friendly as possible, bit there are at least 3 distinct parts:
    1. parsing errors, e.g. a completely invalid JSON/YAML/etc. file - we can't continue from this, we can only show as many details as possible
    2. parsing and type errors for specific fields (e.g. trying to pass a string as an int) - ideally, we should be able to collect all of these errors from all of the sources (CLI, env vars, JSON, etc.) and show them in a single user-friendly list
    3. validation - this is tricky, it's the last step (i.e. we only validate the final consolidated values) and validation logic can spread between multiple fields (e.g. option `X` should be less than or equal to option `Y`)

### Proposed TODO:
0. Figure out a usable Go API (e.g. with initial support for just a few Go types like `string`, `int64` and `bool` that satisfies the criteria :arrow_up: :sweat_smile:
1. Write a PoC with some tests and mock real-life usage examples
2. Iterate and expand on :arrow_up:
3. Support all types (incl. custom types) and multiple sources
4. Polish, set up GitHub Actions CI, etc.
5. Profit

