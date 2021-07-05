package config

import (
	"go.k6.io/croconf"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/executor"
)

type scenariosField struct {
	dest   *lib.ScenarioConfigs
	source *croconf.SourceJSON
	id     string
}

func (sf *scenariosField) Destination() interface{} {
	return sf.dest
}

func (sf *scenariosField) Bindings() []croconf.Binding {
	return []croconf.Binding{
		croconf.NewCallbackBinding(func() error {
			*sf.dest = lib.ScenarioConfigs{
				lib.DefaultScenarioName: executor.NewPerVUIterationsConfig(lib.DefaultScenarioName),
			}
			return nil
		}),
		croconf.NewCallbackBindingFromSource(sf.source, sf.id, func() error {
			jsonVal, ok := sf.source.Lookup(sf.id)
			if !ok {
				return nil
			}

			return sf.dest.UnmarshalJSON(jsonVal)
		}),
	}
}
